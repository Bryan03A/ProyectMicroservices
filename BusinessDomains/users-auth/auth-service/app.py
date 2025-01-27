from werkzeug.security import check_password_hash
from flask import Flask, request, jsonify
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy.orm import Session
import jwt
import datetime
import hashlib
import os

app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = (
    'postgresql://postgres.imfqyzgimtercyyqeqof:1997Guallaba@aws-0-us-west-1.pooler.supabase.com:5432/postgres'
)

app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
app.config['SECRET_KEY'] = 'mysecretkey'

db = SQLAlchemy(app)

# User model
class User(db.Model):
    id = db.Column(db.String(36), primary_key=True)
    username = db.Column(db.String(80), unique=True, nullable=False)
    password = db.Column(db.String(200), nullable=False)
    email = db.Column(db.String(120), unique=True, nullable=False)


# Function to hash the password with the same "salt"
def hash_password(password):
    salt = 'salt'  # Make sure you use the same salt
    dklen = 64 # Damn, that gave problems and it was a matter of putting the 64 in the length
    return hashlib.pbkdf2_hmac('sha256', password.encode('utf-8'), salt.encode('utf-8'), 1000, dklen).hex()

@app.route('/login', methods=['POST'])
def login():
    data = request.get_json()

    # Search for the user by username or email
    user = User.query.filter(
        (User.username == data['username']) | (User.email == data['username'])
    ).first()

    if user:
        # Verify password
        hashed_input_password = hash_password(data['password'])

        if hashed_input_password == user.password:
            # Generate the JWT token
            token = jwt.encode(
                {'user_id': str(user.id), 'exp': datetime.datetime.utcnow() + datetime.timedelta(hours=1)},
                app.config['SECRET_KEY'],
                algorithm='HS256'
            )
            return jsonify({'token': token})

    return jsonify({'message': 'Invalid credentials'}), 401


@app.route('/profile', methods=['GET'])
def profile():
    token = request.headers.get('Authorization')
    if not token:
        return jsonify({"message": "Token is missing"}), 403

    try:
        token = token.split(" ")[1]
        decoded_token = jwt.decode(token, app.config['SECRET_KEY'], algorithms=['HS256'])
        user_id = decoded_token['user_id']

        user = db.session.get(User, user_id)
        if user:
            return jsonify({
                "username": user.username,
                "email": user.email
            })
        else:
            return jsonify({"message": "User not found"}), 404
    except jwt.ExpiredSignatureError:
        return jsonify({"message": "Token expired"}), 401
    except jwt.InvalidTokenError:
        return jsonify({"message": "Invalid token"}), 401

if __name__ == '__main__':
    app.run(debug=True, port=5001)