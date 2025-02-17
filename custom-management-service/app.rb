require 'jwt'
require 'net/http'
require 'uri'
require 'json'
require 'sinatra'
require 'mongo'
require 'pg'
require 'logger'
require 'dotenv/load'
require 'active_record'
require_relative 'price_calculator' # Require the price_calculator.rb file
require 'rack/cors'
require 'bson'  # Make sure to include this library for ObjectId

use Rack::Cors do
  allow do
    origins '*' # You can change '*' to a specific domain if you want to restrict access further
    resource '*', headers: :any, methods: [:get, :post, :put, :delete, :options]
  end
end

set :port, 5012

# MongoDB connection setup
client = Mongo::Client.new(ENV['MONGO_URI'])
db = client.database
models_collection = db[:models]

# Setup database connection using ActiveRecord
ActiveRecord::Base.establish_connection(
  ENV['POSTGRES_URI'] + '?prepared_statements=false'
)

# Define the 'custom' table if it does not exist
class Custom < ActiveRecord::Base
  validates :model_id, :user_id, :created_by, :cost_initial, :cost_final, presence: true
end

# Verify if the 'custom_models' table exists, if not, create it
conn = ActiveRecord::Base.connection.raw_connection
conn.exec <<-SQL
  CREATE TABLE IF NOT EXISTS customs (
    id SERIAL PRIMARY KEY,
    model_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_by TEXT NOT NULL,
    custom_params JSONB NOT NULL,
    cost_initial DECIMAL NOT NULL,
    cost_final DECIMAL NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );
SQL

# Route to get models from the catalog in MongoDB
get '/models' do
  models = models_collection.find.map { |model| model }
  models.to_json
end

# Function to verify the token with the auth-service
def verify_token(token)
  uri = URI.parse("http://localhost:5001/profile")
  req = Net::HTTP::Get.new(uri)
  req['Authorization'] = "Bearer #{token}"

  res = Net::HTTP.start(uri.hostname, uri.port) { |http| http.request(req) }
  if res.code.to_i == 200
    user_info = JSON.parse(res.body)
    return user_info['username'], user_info['email'], user_info['user_id']
  else
    raise 'Invalid or expired token'
  end
end

post '/customize-model' do
  begin
    # Get the token from the header
    token = request.env["HTTP_AUTHORIZATION"]&.split(" ")&.last
    raise 'Token missing' unless token

    # Verify the token with auth-service
    username, email, user_id = verify_token(token)
    puts "✅ Authenticated user: ID=#{user_id}" # Log the user ID

    # Get the model data and custom parameters
    data = JSON.parse(request.body.read)
    model_id = data['model_id']
    custom_params = data['custom_params']

    # Get the model from MongoDB using the model_id
    model = models_collection.find(_id: BSON::ObjectId(model_id)).first
    return status 404 if model.nil?

    # Get the model price from MongoDB
    price = model['price'].to_f
    puts "Initial model price: #{price}"

    # Get the 'created_by' field from the model (name of the user who created it)
    created_by = model['created_by']
    puts "Model created by: #{created_by}"

    # Use the PriceCalculator class to calculate the final price
    cost_final = PriceCalculator.calculate(price, custom_params)
    puts "Calculated final price: #{cost_final}"

    # Save the customization in the database
    custom_model = Custom.create(
      model_id: model_id,
      user_id: user_id,
      created_by: created_by,  # Add 'created_by' here
      custom_params: custom_params.to_json,
      cost_initial: price,
      cost_final: cost_final
    )

    status 200
    { message: 'Model customized and saved successfully', custom_model: custom_model.to_json }.to_json
  rescue => e
    status 400
    { message: "Error: #{e.message}" }.to_json
  end
end