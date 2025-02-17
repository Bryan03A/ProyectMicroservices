from flask import Flask, render_template, request, jsonify
from flask_socketio import SocketIO, emit, join_room
from flask_cors import CORS
import psycopg2
import os
from pymongo import MongoClient
from bson.objectid import ObjectId
import datetime  # Necesario para los timestamps

app = Flask(__name__)
socketio = SocketIO(app, cors_allowed_origins="*")  # Solo una instancia de SocketIO

# Configura CORS para permitir solicitudes desde cualquier origen (o desde un origen específico)
CORS(app, resources={r"/*": {"origins": "*"}})

# Conexión a MongoDB
mongo_uri = "mongodb+srv://MicroserviceDev:1997999@cluster0.hdqpd.mongodb.net/CatalogServiceDB?retryWrites=true&w=majority"
client = MongoClient(mongo_uri)
db = client.get_database()
chats_collection = db.chats

# Conexión a PostgreSQL
postgres_uri = "postgresql://postgres.imfqyzgimtercyyqeqof:1997Guallaba@aws-0-us-west-1.pooler.supabase.com:6543/postgres"
conn = psycopg2.connect(postgres_uri)
cur = conn.cursor()

# Crear tabla si no existe
cur.execute("""
    CREATE TABLE IF NOT EXISTS chat_sessions (
        id SERIAL PRIMARY KEY,
        id_cliente VARCHAR(255),
        created_by VARCHAR(255),
        id_modelo VARCHAR(255),
        id_chat VARCHAR(255)
    );
""")
conn.commit()

# Ruta para la página principal (chat)
@app.route('/')
def index():
    return render_template('index.html')

# Ruta para crear un chat
@app.route('/start_chat', methods=['POST'])
def start_chat():
    data = request.get_json()
    id_cliente = data['id_cliente']
    created_by = data['created_by']
    id_modelo = data['id_modelo']
    
    # Verificar si ya existe un chat entre este cliente y vendedor
    existing_chat = chats_collection.find_one({
        '$or': [
            {'id_cliente': id_cliente, 'created_by': created_by},
            {'id_cliente': created_by, 'created_by': id_cliente}
        ],
        'id_modelo': id_modelo
    })
    
    if existing_chat:
        return jsonify({"status": "success", "id_chat": str(existing_chat["_id"])})
    
    # Si no existe, crear un nuevo chat
    chat_data = {
        "id_cliente": id_cliente,
        "created_by": created_by,
        "id_modelo": id_modelo,
        "messages": []
    }
    
    new_chat = chats_collection.insert_one(chat_data)
    
    # Guardar el chat en PostgreSQL
    cur.execute("""
        INSERT INTO chat_sessions (id_cliente, created_by, id_modelo, id_chat)
        VALUES (%s, %s, %s, %s);
    """, (id_cliente, created_by, id_modelo, str(new_chat.inserted_id)))
    conn.commit()

    return jsonify({"status": "success", "id_chat": str(new_chat.inserted_id)})

# Manejo de eventos de Socket.IO
@socketio.on('send_message')
def handle_message(data):
    # Guardar mensaje en MongoDB
    chat_id = data['id_chat']
    message = data['message']
    sender = data['sender']  # cliente o vendedor

    chat = chats_collection.find_one({'_id': ObjectId(chat_id)})
    if chat:
        # Actualizar los mensajes del chat en MongoDB
        chat['messages'].append({
            'sender': sender,
            'message': message,
            'timestamp': str(datetime.datetime.now())
        })
        chats_collection.update_one({'_id': ObjectId(chat_id)}, {'$set': {'messages': chat['messages']}})
    
    # Emitir el mensaje a todos los participantes del chat
    emit('receive_message', data, room=chat_id)

# Ruta para obtener el historial de mensajes de un chat
@app.route('/api/chats/<chat_id>', methods=['GET'])
def get_chat_history(chat_id):
    # Buscar el historial de chat desde MongoDB
    chat_history = chats_collection.find_one({'_id': ObjectId(chat_id)})
    if chat_history:
        return jsonify(chat_history['messages']), 200
    return jsonify({'message': 'Chat not found'}), 404

# WebSocket para manejo de chat en tiempo real
@socketio.on('connect')
def handle_connect():
    print("Nuevo cliente conectado")

@socketio.on('disconnect')
def handle_disconnect():
    print("Cliente desconectado")

@socketio.on('join_chat')
def handle_join_chat(data):
    chat_id = data.get('chatId')
    user_name = data.get('userName')
    # Unirse al chat
    print(f"{user_name} se unió al chat: {chat_id}")
    join_room(chat_id)
    
    # Enviar historial de mensajes al cliente
    chat_history = chats_collection.find_one({'_id': ObjectId(chat_id)})
    if chat_history:
        emit('load_history', chat_history['messages'], room=chat_id)

if __name__ == '__main__':
    socketio.run(app, host='0.0.0.0', port=5010)