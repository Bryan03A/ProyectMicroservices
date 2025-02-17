from werkzeug.security import check_password_hash
from flask import Flask, request, jsonify
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy import text
import jwt
import datetime
import hashlib
from flask_cors import CORS
from dotenv import load_dotenv
import os

# Load environment variables from the .env file
load_dotenv()

# Initialize Flask app
app = Flask(__name__)

# Enable CORS to allow requests from specific frontend domain (localhost:8080 in this case)
CORS(app, origins=["http://52.91.86.137:8080"], supports_credentials=True)

# Configuring the connection to PostgreSQL (AWS RDS)
app.config['SQLALCHEMY_DATABASE_URI'] = os.getenv('DATABASE_URL')  # Use DATABASE_URL from environment
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = os.getenv('SQLALCHEMY_TRACK_MODIFICATIONS', 'False') == 'True'
app.config['SECRET_KEY'] = os.getenv('SECRET_KEY')  # Secret key for JWT encoding and decoding

# Initialize the SQLAlchemy object to interact with the database
db = SQLAlchemy(app)

# Test the database connection with a simple query
@app.route('/test-db', methods=['GET'])
def test_db():
    try:
        result = db.session.execute(text("SELECT NOW()")).fetchone()  # Executes a SQL query to check DB status
        return jsonify({"message": "Conexión exitosa", "timestamp": str(result[0])})
    except Exception as e:
        return jsonify({"error": str(e)}), 500  # Return an error message if connection fails

# User model for interacting with the database
class User(db.Model):
    id = db.Column(db.String(36), primary_key=True)  # User ID (UUID)
    username = db.Column(db.String(80), unique=True, nullable=False)  # Username (must be unique)
    password = db.Column(db.String(200), nullable=False)  # Password (hashed)
    email = db.Column(db.String(120), unique=True, nullable=False)  # Email (must be unique)

# Function to hash passwords securely using PBKDF2 algorithm
def hash_password(password):
    salt = 'salt'  # Static salt for simplicity, should ideally be random
    dklen = 64  # Length of the derived key
    return hashlib.pbkdf2_hmac('sha256', password.encode('utf-8'), salt.encode('utf-8'), 1000, dklen).hex()

# Login endpoint to authenticate users
@app.route('/login', methods=['POST'])
def login():
    data = request.get_json()  # Extract JSON data from the request

    # Search for the user by username or email
    user = User.query.filter(
        (User.username == data['username']) | (User.email == data['username'])
    ).first()

    if user:
        # Hash the entered password and compare it with the stored password
        hashed_input_password = hash_password(data['password'])

        if hashed_input_password == user.password:
            # If passwords match, create a JWT token
            token = jwt.encode(
                {
                    'user_id': str(user.id),  # User's ID in the token payload
                    'username': user.username,  # Username in the token payload
                    'exp': datetime.datetime.utcnow() + datetime.timedelta(hours=1)  # Expiry time (1 hour)
                },
                app.config['SECRET_KEY'],  # Secret key used for signing the token
                algorithm='HS256'  # Algorithm for encoding the token
            )
            return jsonify({'token': token})  # Return the JWT token as the response

    return jsonify({'message': 'Invalid credentials'}), 401  # If credentials are invalid, return error

# Profile endpoint to retrieve user profile based on the JWT token
@app.route('/profile', methods=['GET'])
def profile():
    token = request.headers.get('Authorization')  # Get the Authorization header containing the token
    if not token:
        return jsonify({"message": "Token is missing"}), 403  # If token is not provided, return error

    try:
        # Extract the token (bearer token format)
        token = token.split(" ")[1]
        # Decode the JWT token using the secret key
        decoded_token = jwt.decode(token, app.config['SECRET_KEY'], algorithms=['HS256'])
        user_id = decoded_token['user_id']  # Extract user_id from the decoded token

        # Fetch user from the database using the decoded user_id
        user = db.session.get(User, user_id)
        if user:
            # Return user profile information (username and email)
            return jsonify({
                "username": user.username,
                "email": user.email,
                "user_id": user_id
            })
        else:
            return jsonify({"message": "User not found"}), 404  # If user is not found, return error

    except jwt.ExpiredSignatureError:
        return jsonify({"message": "Token expired"}), 401  # If the token is expired
    except jwt.InvalidTokenError:
        return jsonify({"message": "Invalid token"}), 401  # If the token is invalid

# Run the Flask app
if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5001)