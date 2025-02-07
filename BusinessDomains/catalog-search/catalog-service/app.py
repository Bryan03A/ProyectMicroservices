from flask import Flask, jsonify, request
from pymongo import MongoClient
import os

app = Flask(__name__)

# MongoDB Configuration
MONGO_URI = "mongodb+srv://MicroserviceDev:1997999@cluster0.hdqpd.mongodb.net/CatalogServiceDB?retryWrites=true&w=majority"
client = MongoClient(MONGO_URI)
db = client["CatalogServiceDB"]
models_collection = db["models"]  # Collection to store 3D models

# Test route
@app.route("/")
def home():
    return jsonify({"message": "Welcome to the Catalog Service!"})

# Route to add a new 3D model
@app.route("/models", methods=["POST"])
def add_model():
    try:
        model_data = request.json
        result = models_collection.insert_one(model_data)
        return jsonify({"message": "Model added successfully", "model_id": str(result.inserted_id)}), 201
    except Exception as e:
        return jsonify({"error": str(e)}), 500

# Route to get all 3D models
@app.route("/models", methods=["GET"])
def get_models():
    try:
        models = list(models_collection.find({}, {"_id": 0}))  # Exclude the _id field
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
        result = models_collection.delete_one({"name": model_name})
        if result.deleted_count > 0:
            return jsonify({"message": "Model deleted successfully"}), 200
        else:
            return jsonify({"error": "Model not found"}), 404
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    try:
        # Check MongoDB connection on startup
        client.admin.command("ping")
        print("Connected to MongoDB Atlas")
    except Exception as e:
        print(f"Error connecting to MongoDB: {e}")
    app.run(debug=True, port=5003)