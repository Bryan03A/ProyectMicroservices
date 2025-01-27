const express = require('express');
const { Pool } = require('pg');
const { v4: uuidv4 } = require('uuid');
const crypto = require('crypto');  // We use crypto to create hashes
const app = express();
const port = 5002;

// Connecting to PostgreSQL (Supabase)
const pool = new Pool({
    connectionString: 'postgresql://postgres.imfqyzgimtercyyqeqof:1997Guallaba@aws-0-us-west-1.pooler.supabase.com:5432/postgres?pool_mode=session',
});

app.use(express.json()); // For the server to accept JSON

// Create table if it does not exist
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

// Function to make the hash with pbkdf2:sha256
const hashPassword = (password) => {
  return new Promise((resolve, reject) => {
      const salt = 'salt'; // Make sure you use the same salt
      const iterations = 1000;
      const keyLength = 64; // I see you damn

      crypto.pbkdf2(password, salt, iterations, keyLength, 'sha256', (err, derivedKey) => {
          if (err) reject(err);
          const hashedPassword = derivedKey.toString('hex');
          resolve(hashedPassword); // Returns the hash in hexadecimal format
      });
  });
};

// Route to register a new user
app.post('/register', async (req, res) => {
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

// Start server
app.listen(port, async () => {
    await createTable();
    console.log(`User-service is running on port ${port}`);
});