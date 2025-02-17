# model_service.py

from flask import Flask, jsonify, request
from pymongo import MongoClient
import requests
from dotenv import load_dotenv
import os

load_dotenv()

app = Flask(__name__)

# MongoDB Configuration
MONGO_URI = os.getenv("MONGO_URI")
client = MongoClient(MONGO_URI)
db = client["CatalogServiceDB"]
models_collection = db["models"]

# Elasticsearch URL
ELASTICSEARCH_URL = "http://localhost:9200"

# Helper function to index a model in Elasticsearch
def index_model_in_elasticsearch(model_data, model_id):
    try:
        response = requests.post(f"{ELASTICSEARCH_URL}/models/_doc/{model_id}", json=model_data)
        if response.status_code in [200, 201]:
            print("Model indexed in Elasticsearch.")
        else:
            print("Error indexing model in Elasticsearch.")
    except requests.exceptions.RequestException as e:
        print(f"Error: {e}")

# Route to add a new model
@app.route("/models", methods=["POST"])
def add_model():
    try:
        model_data = request.json
        result = models_collection.insert_one(model_data)
        model_id = str(result.inserted_id)
        
        # Index model in Elasticsearch
        index_model_in_elasticsearch(model_data, model_id)

        return jsonify({"message": "Model added successfully", "model_id": model_id}), 201
    except Exception as e:
        return jsonify({"error": str(e)}), 500

# Route to get all models
@app.route("/models", methods=["GET"])
def get_models():
    try:
        models = list(models_collection.find({}))
        for model in models:
            model["_id"] = str(model["_id"])  # Convert ObjectId to string
        return jsonify({"models": models}), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=5012)
