from flask import Flask, jsonify, request
from pymongo import MongoClient
import jwt
from functools import wraps
from flask_cors import CORS
from dotenv import load_dotenv
import os

load_dotenv()

app = Flask(__name__)
CORS(app)

# Enable CORS to allow requests from localhost:8080
CORS(app, origins=["http://3.82.92.84:8080"], supports_credentials=True)

# MongoDB Configuration
MONGO_URI = os.getenv("MONGO_URI")
client = MongoClient(MONGO_URI)
db = client["CatalogServiceDB"]
models_collection = db["models"]  # Collection to store 3D models

# Secret key for JWT
SECRET_KEY = os.getenv("SECRET_KEY")

# Function to decode JWT token and get user_id
def get_user_info_from_token():
    token = request.headers.get('Authorization')
    if not token:
        return None, None
    try:
        token = token.split(" ")[1]  # Asume que el token está precedido por "Bearer"
        decoded = jwt.decode(token, SECRET_KEY, algorithms=["HS256"], options={"verify_exp": False})
        user_id = decoded.get('user_id')
        user_name = decoded.get('username')  # Extrae el username del token
        print(f"Decoded token: user_id={user_id}, username={user_name}")  # Agregar esta línea para ver los valores
        return user_id, user_name
    except jwt.PyJWTError as e:
        print(f"JWT Error: {str(e)}")
        return None, None
    except Exception as e:
        print(f"Error: {str(e)}")
        return None, None

# Middleware to check if the current user is the owner of the model
def check_model_owner(model_id):
    user_id, user_name = get_user_info_from_token()
    if not user_name:  # Aquí debería verificar el nombre de usuario, no el user_id
        return None
    model = models_collection.find_one({"_id": model_id})
    if model and model.get('created_by') == user_name:  # Compara user_name con created_by
        return user_name
    return None

# Test route
@app.route("/")
def home():
    return jsonify({"message": "Welcome to the Catalog Service!"})

@app.route("/models", methods=["POST"])
def add_model():
    try:
        model_data = request.json
        print("Received model data:", model_data)

        user_id, user_name = get_user_info_from_token()
        if not user_id:
            return jsonify({"error": "Unauthorized, no token provided"}), 401
        
        if not user_name:
            return jsonify({"error": "User name not found in token"}), 401
        
        model_data['created_by'] = user_name  # Usamos el nombre del usuario
        result = models_collection.insert_one(model_data)
        return jsonify({"message": "Model added successfully", "model_id": str(result.inserted_id)}), 201
    except Exception as e:
        print(f"Error: {str(e)}")
        return jsonify({"error": f"An error occurred: {str(e)}"}), 500

# Route to get all 3D models
@app.route("/models", methods=["GET"])
def get_models():
    try:
        models = list(models_collection.find({}, {"_id": 0}))  # Exclude the _id field
        return jsonify({"models": models}), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500

# Route to get all 3D models by user ID
@app.route("/models/user/<string:user_id>", methods=["GET"])
def get_models_by_user(user_id):
    try:
        models = list(models_collection.find({"created_by": user_id}, {"_id": 0}))  # Filter by created_by
        return jsonify({"models": models}), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500

# Route to get a specific 3D model by its name
@app.route("/models/<string:model_name>", methods=["GET"])
def get_model(model_name):
    try:
        model = models_collection.find_one({"name": model_name}, {"_id": 0})
        if model:
            return jsonify({"model": model}), 200
        else:
            return jsonify({"error": "Model not found"}), 404
    except Exception as e:
        return jsonify({"error": str(e)}), 500

# Route to delete a 3D model by its name
@app.route("/models/<string:model_name>", methods=["DELETE"])
def delete_model(model_name):
    try:
        model = models_collection.find_one({"name": model_name})
        if not model:
            return jsonify({"error": "Model not found"}), 404
        user_name = check_model_owner(model['_id'])  # Pasar _id, no user_id
        if not user_name:
            return jsonify({"error": "Unauthorized"}), 403  # Only owner can delete
        result = models_collection.delete_one({"name": model_name})
        if result.deleted_count > 0:
            return jsonify({"message": "Model deleted successfully"}), 200
        else:
            return jsonify({"error": "Error deleting model"}), 500
    except Exception as e:
        return jsonify({"error": str(e)}), 500

# Route to update a 3D model by its name
@app.route("/models/<string:model_name>", methods=["PUT"])
def update_model(model_name):
    try:
        model = models_collection.find_one({"name": model_name})
        if not model:
            return jsonify({"error": "Model not found"}), 404
        if not check_model_owner(model['_id']):
            return jsonify({"error": "Unauthorized"}), 403  # Only owner can update
        updated_data = request.json
        result = models_collection.update_one({"name": model_name}, {"$set": updated_data})
        if result.modified_count > 0:
            return jsonify({"message": "Model updated successfully"}), 200
        else:
            return jsonify({"error": "Error updating model"}), 500
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    try:
        # Check MongoDB connection on startup
        client.admin.command("ping")
        print("Connected to MongoDB Atlas")
    except Exception as e:
        print(f"Error connecting to MongoDB: {e}")
    app.run(debug=True, host='0.0.0.0', port=5003)