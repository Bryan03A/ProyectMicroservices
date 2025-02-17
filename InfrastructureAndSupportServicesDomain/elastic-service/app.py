# elastic_service.py

from flask import Flask, jsonify, request
import requests
from dotenv import load_dotenv
import os

load_dotenv()

app = Flask(__name__)

# Elasticsearch URL
ELASTICSEARCH_URL = "http://localhost:9200"

# Route to index a model in Elasticsearch
@app.route("/elasticsearch/models", methods=["POST"])
def index_model():
    try:
        model_data = request.json
        model_id = model_data.get("model_id")
        response = requests.post(f"{ELASTICSEARCH_URL}/models/_doc/{model_id}", json=model_data)
        
        if response.status_code == 201:
            return jsonify({"message": "Model indexed successfully"}), 201
        else:
            return jsonify({"error": "Error indexing model"}), 500
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500

# Route to delete a model from Elasticsearch
@app.route("/elasticsearch/models/<model_id>", methods=["DELETE"])
def delete_model(model_id):
    try:
        response = requests.delete(f"{ELASTICSEARCH_URL}/models/_doc/{model_id}")
        
        if response.status_code == 200:
            return jsonify({"message": "Model deleted successfully"}), 200
        else:
            return jsonify({"error": "Error deleting model"}), 500
    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=5016)