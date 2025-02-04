require('dotenv').config();

const express = require('express');
const { Pool } = require('pg');
const { v4: uuidv4 } = require('uuid');
const crypto = require('crypto');  
const cors = require('cors');  // Importar CORS

const app = express();
const port = 5002;

// Middleware para procesar JSON y habilitar CORS
app.use(express.json());  // <-- Permite recibir JSON en req.body
app.use(express.urlencoded({ extended: true }));  // <-- Para datos de formularios
app.use(cors());  // <-- Permite solicitudes desde otros dominios

// Configurar la conexión a PostgreSQL en AWS RDS
const pool = new Pool({
  host: process.env.POSTGRESQL_HOST,
  port: process.env.POSTGRESQL_PORT,
  database: process.env.POSTGRESQL_DATABASE,
  user: process.env.POSTGRESQL_USER,
  password: process.env.POSTGRESQL_PASSWORD,
  ssl: {
    rejectUnauthorized: false
  }
});

// Probar la conexión
async function testConnection() {
  try {
    const client = await pool.connect();
    const result = await client.query('SELECT NOW()');
    console.log('Conectado a PostgreSQL en AWS RDS:', result.rows[0]);
    client.release();
  } catch (error) {
    console.error('Error de conexión:', error);
  }
}

testConnection();

// Crear tabla si no existe
async function createTable() {
  const query = `
      CREATE TABLE IF NOT EXISTS "user" (
          id UUID PRIMARY KEY,
          username VARCHAR(80) UNIQUE NOT NULL,
          password VARCHAR(200) NOT NULL,
          first_name VARCHAR(100) NOT NULL,
          last_name VARCHAR(100) NOT NULL,
          dni VARCHAR(20) UNIQUE NOT NULL,
          email VARCHAR(120) UNIQUE NOT NULL,
          city VARCHAR(100) NOT NULL
      );
  `;
  try {
      await pool.query(query);
      console.log('Table "user" created or already existed');
  } catch (error) {
      console.error('Error creating table:', error);
  }
}

// Función para hashear contraseñas
const hashPassword = (password) => {
  return new Promise((resolve, reject) => {
      const salt = 'salt';
      const iterations = 1000;
      const keyLength = 64;

      crypto.pbkdf2(password, salt, iterations, keyLength, 'sha256', (err, derivedKey) => {
          if (err) reject(err);
          resolve(derivedKey.toString('hex'));
      });
  });
};

// **RUTA PARA REGISTRAR UN USUARIO**
app.post('/register', async (req, res) => {
  console.log('Body recibido:', req.body);  // <-- Verifica que req.body no sea undefined

  const { username, password, first_name, last_name, dni, email, city } = req.body;

  if (!username || !password || !first_name || !last_name || !dni || !email || !city) {
      return res.status(400).json({ message: 'All fields are required' });
  }

  try {
      const hashedPassword = await hashPassword(password);

      const userId = uuidv4();
      const query = `
          INSERT INTO "user" (id, username, password, first_name, last_name, dni, email, city)
          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
      `;
      await pool.query(query, [userId, username, hashedPassword, first_name, last_name, dni, email, city]);
      res.status(201).json({ message: 'User registered successfully' });
  } catch (error) {
      console.error('Error registering user:', error);
      res.status(500).json({ message: 'Error registering user' });
  }
});

// **RUTA PARA OBTENER UN USUARIO POR ID**
app.get('/user/:id', async (req, res) => {
    const { id } = req.params;
  
    try {
      const query = `
        SELECT id, username, first_name, last_name, dni, email, city
        FROM "user"
        WHERE id = $1
      `;
      
      const result = await pool.query(query, [id]);
  
      if (result.rows.length > 0) {
        res.status(200).json(result.rows[0]);
      } else {
        res.status(404).json({ message: 'User not found' });
      }
    } catch (error) {
      console.error('Error fetching user information:', error);
      res.status(500).json({ message: 'Error fetching user information' });
    }
});

// Iniciar servidor
app.listen(port, '0.0.0.0', async () => {
    await createTable();
    console.log(`User-service is running on port ${port}`);
});