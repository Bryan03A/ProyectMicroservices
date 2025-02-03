require('dotenv').config();

const express = require('express');
const mongoose = require('mongoose');
const sessionRoutes = require('./routes/sessions');
const authenticateToken = require('./middlewares/authMiddleware');
const cors = require('cors');

// MongoDB connection string
const MONGO_URI = process.env.MONGO_URI;
const PORT = process.env.PORT || 5004;

if (!MONGO_URI) {
    console.error("❌ ERROR: MONGO_URI no está definido. Verifica tu archivo .env.");
    process.exit(1); // Detiene la ejecución si falta la conexión
}

const app = express();


app.use(cors({
    origin: 'http://localhost:8080',  // Permite solicitudes solo desde tu frontend
    credentials: true  // Permite enviar cookies y headers de autenticación
}));

// Middleware
app.use(express.json());
app.use('/api/sessions', sessionRoutes);

// Connect to MongoDB
mongoose.set('strictQuery', false);
mongoose.connect(MONGO_URI, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
})
.then(() => console.log('Connected to MongoDB'))
.catch((err) => console.error('Error connecting to MongoDB:', err));

// Start server
app.listen(PORT, () => {
    console.log(`Session Service running on port ${PORT}`);
});