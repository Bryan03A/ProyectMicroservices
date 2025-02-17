# auth_service.py

from flask import Flask, jsonify, request
import jwt
import os
from datetime import datetime, timedelta
from dotenv import load_dotenv

load_dotenv()

app = Flask(__name__)

# Secret key for JWT
SECRET_KEY = os.getenv("SECRET_KEY")

# Route to authenticate and generate JWT
@app.route("/auth/login", methods=["POST"])
def login():
    try:
        data = request.json
        username = data.get("username")
        password = data.get("password")

        # In a real application, validate the username and password from a DB
        if username and password:  # Simulating valid authentication
            payload = {
                "user_id": username,  # In a real scenario, this would be a user ID from the DB
                "username": username,
                "exp": datetime.utcnow() + timedelta(hours=1)
            }
            token = jwt.encode(payload, SECRET_KEY, algorithm="HS256")
            return jsonify({"token": token}), 200
        else:
            return jsonify({"error": "Invalid credentials"}), 401
    except Exception as e:
        return jsonify({"error": str(e)}), 500

# Route to validate JWT
@app.route("/auth/validate", methods=["POST"])
def validate_token():
    token = request.headers.get("Authorization")
    if not token:
        return jsonify({"error": "No token provided"}), 401
    
    try:
        token = token.split(" ")[1]  # "Bearer <token>"
        decoded = jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
        return jsonify({"message": "Token is valid", "user_id": decoded["user_id"]}), 200
    except jwt.ExpiredSignatureError:
        return jsonify({"error": "Token has expired"}), 401
    except jwt.InvalidTokenError:
        return jsonify({"error": "Invalid token"}), 401

if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=5011)