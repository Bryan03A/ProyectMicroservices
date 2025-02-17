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

# Check if the 'custom_models' table exists, if not, create it
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

# Route to get models from the catalog from MongoDB
get '/models' do
  models = models_collection.find.map { |model| model }
  models.to_json
end

# Function to verify the token with the auth-service
def verify_token(token)
  uri = URI.parse("http://52.91.86.137:5001/profile")
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
    puts "✅ User authenticated: ID=#{user_id}" # Log the user ID

    # Get the model data and custom parameters
    data = JSON.parse(request.body.read)
    model_id = data['model_id']
    custom_params = data['custom_params']

    # Get the model from MongoDB using the model_id
    model = models_collection.find(_id: BSON::ObjectId(model_id)).first
    return status 404 if model.nil?

    # Get the price of the model from MongoDB
    price = model['price'].to_f
    puts "Initial price of the model: #{price}"

    # Get the 'created_by' field from the model (the username of the creator)
    created_by = model['created_by']
    puts "Model created by: #{created_by}"

    # Use the PriceCalculator class to calculate the final price
    cost_final = PriceCalculator.calculate(price, custom_params)
    puts "Calculated final price: #{cost_final}"

    # Save the customization to the database
    custom_model = Custom.create(
      model_id: model_id,
      user_id: user_id,
      created_by: created_by,  # Add the 'created_by' here
      custom_params: custom_params.to_json,
      cost_initial: price,
      cost_final: cost_final
    )

    status 200
    { message: 'Model customized and saved successfully', custom_model: custom_model }.to_json
  rescue => e
    status 400
    { message: e.message }.to_json
  end
end

# Endpoint to update a customized model
put '/customize-model/:id' do
  begin
    # Get the token from the header
    token = request.env["HTTP_AUTHORIZATION"]&.split(" ")&.last
    raise 'Token missing' unless token

    # Verify the token with auth-service
    username, email, user_id = verify_token(token)

    # Find the custom model
    custom_model = Custom.find_by(id: params[:id])
    raise 'Custom model not found' unless custom_model

    # Check if the user is authorized to edit this model
    raise 'Unauthorized' unless custom_model.user_id == user_id

    # Get new customization data
    data = JSON.parse(request.body.read)
    custom_params = data['custom_params'] # Get the new custom parameters

    # Update the custom_params field in the correct JSON format
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

    # Verify the token with auth-service
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

# Endpoint to list all customized models for a user
get '/customize-models' do
  begin
    # Get the token from the header
    token = request.env["HTTP_AUTHORIZATION"]&.split(" ")&.last
    raise 'Token missing' unless token

    # Verify the token with auth-service
    username, email, user_id = verify_token(token)

    # Fetch all custom models for the user
    custom_models = Custom.where(user_id: user_id)

    status 200
    custom_models.to_json
  rescue => e
    status 400
    { message: e.message }.to_json
  end
end

# Endpoint to get a customized model by ID
get '/customize-model/:id' do
  begin
    # Get the token from the header
    token = request.env["HTTP_AUTHORIZATION"]&.split(" ")&.last
    raise 'Token missing' unless token

    # Verify the token with auth-service
    username, email, user_id = verify_token(token)

    # Find the custom model by ID
    custom_model = Custom.find_by(id: params[:id])
    raise 'Custom model not found' unless custom_model

    # Check if the user is authorized to view this model
    raise 'Unauthorized' unless custom_model.user_id == user_id

    # Fetch the original model from MongoDB using the model_id
    model = models_collection.find(_id: BSON::ObjectId(custom_model.model_id)).first
    raise 'Original model not found' unless model

    # Return the customized model details along with the original model
    status 200
    {
      custom_model: custom_model,
      original_model: model
    }.to_json
  rescue => e
    status 400
    { message: e.message }.to_json
  end
end

set :port, 5007