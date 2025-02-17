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

use Rack::Cors do
  allow do
    origins '*' # You can change '*' to a specific domain if you want to restrict access further
    resource '*', headers: :any, methods: [:get, :post, :put, :delete, :options]
  end
end

set :port, 5014

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

# Route to fetch models from the catalog (MongoDB)
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

# Endpoint to update a customized model
put '/customize-model/:id' do
  begin
    # Get the token from the header
    token = request.env["HTTP_AUTHORIZATION"]&.split(" ")&.last
    raise 'Token missing' unless token

    # Verify the token with the auth-service
    username, email, user_id = verify_token(token)

    # Find the custom model
    custom_model = Custom.find_by(id: params[:id])
    raise 'Custom model not found' unless custom_model

    # Check if the user is authorized to edit this model
    raise 'Unauthorized' unless custom_model.user_id == user_id

    # Get new customization data
    data = JSON.parse(request.body.read)
    custom_params = data['custom_params'] # Get the new customization parameters

    # Update the custom_params field to JSON format correctly
    custom_model.update(
      custom_params: custom_params.to_json,
      cost_initial: custom_model.cost_initial, # Keep the initial price
      cost_final: PriceCalculator.calculate(custom_model.cost_initial, custom_params) # Recalculate the final price
    )

    status 200
    { message: 'Model customized successfully', custom_model: custom_model }.to_json
  rescue => e
    status 400
    { message: e.message }.to_json
  end
end

# Endpoint to delete a customized model
delete '/customize-model/:id' do
  begin
    # Get the token from the header
    token = request.env["HTTP_AUTHORIZATION"]&.split(" ")&.last
    raise 'Token missing' unless token

    # Verify the token with the auth-service
    username, email, user_id = verify_token(token)

    # Find the custom model
    custom_model = Custom.find_by(id: params[:id])
    raise 'Custom model not found' unless custom_model

    # Check if the user is authorized to delete this model
    raise 'Unauthorized' unless custom_model.user_id == user_id

    # Delete the model
    custom_model.destroy

    status 200
    { message: 'Custom model deleted successfully' }.to_json
  rescue => e
    status 400
    { message: e.message }.to_json
  end
end