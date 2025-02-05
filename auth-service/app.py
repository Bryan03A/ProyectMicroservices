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

load_dotenv()

app = Flask(__name__)

# Enable CORS to allow requests from localhost:8080
CORS(app, origins=["http://54.214.66.77:8080"], supports_credentials=True)

# Configurar la conexión a PostgreSQL (AWS RDS)
app.config['SQLALCHEMY_DATABASE_URI'] = os.getenv('DATABASE_URL')
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = os.getenv('SQLALCHEMY_TRACK_MODIFICATIONS', 'False') == 'True'
app.config['SECRET_KEY'] = os.getenv('SECRET_KEY')

db = SQLAlchemy(app)

# Ruta para probar la conexión con una consulta a la base de datos
@app.route('/test-db', methods=['GET'])
def test_db():
    try:
        result = db.session.execute(text("SELECT NOW()")).fetchone()  # Se usa text() aquí
        return jsonify({"message": "Conexión exitosa", "timestamp": str(result[0])})
    except Exception as e:
        return jsonify({"error": str(e)}), 500

# User model
class User(db.Model):
    id = db.Column(db.String(36), primary_key=True)
    username = db.Column(db.String(80), unique=True, nullable=False)
    password = db.Column(db.String(200), nullable=False)
    email = db.Column(db.String(120), unique=True, nullable=False)

# Function to hash the password
def hash_password(password):
    salt = 'salt'  # Make sure to use the same salt
    dklen = 64  # Length of the derived key
    return hashlib.pbkdf2_hmac('sha256', password.encode('utf-8'), salt.encode('utf-8'), 1000, dklen).hex()

@app.route('/login', methods=['POST'])
def login():
    data = request.get_json()

    # Buscar al usuario por username o email
    user = User.query.filter(
        (User.username == data['username']) | (User.email == data['username'])
    ).first()

    if user:
        # Verificar la contraseña
        hashed_input_password = hash_password(data['password'])

        if hashed_input_password == user.password:
            # Generar el token JWT
            token = jwt.encode(
                {
                    'user_id': str(user.id),
                    'username': user.username,
                    'exp': datetime.datetime.utcnow() + datetime.timedelta(hours=1)
                },
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
    app.run(debug=True, host='0.0.0.0', port=5001)
