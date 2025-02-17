const express = require('express');
const path = require('path');
const axios = require('axios');
const cors = require('cors');

const app = express();
const port = 5018;

// CORS Configuration
app.use(cors({
    origin: '*',  // This allows requests from port 8080
    methods: ['GET', 'POST', 'PUT', 'DELETE'],  // You can adjust the methods as needed
    allowedHeaders: ['Content-Type', 'Authorization'],  // Allowed headers
}));

// Importing routes
const userRoutes = require('./routes/user-service');
const authRoutes = require('./routes/auth-service');
const sessionRoutes = require('./routes/session-service');
const catalogRoutes = require('./routes/catalog-service');
const searchService = require("./routes/search-service");

// Serving static files (such as index.html)
app.use(express.static(path.join(__dirname, 'public')));

// Middleware for parsing JSON
app.use(express.json());

// Use the user routes
app.use(userRoutes);

// Use the auth routes
app.use(authRoutes);

app.use(sessionRoutes);

// Use the catalog routes
app.use(catalogRoutes);

// Ruta para manejar las búsquedas
app.use('/search', require('./routes/search-service'));

// Ruta para la página de chat
app.get('/chat', (req, res) => {
    res.sendFile(path.join(__dirname, 'public', 'chat.html'));
});

// Start server
app.listen(port, '0.0.0.0', () => {
    console.log(`Graphical interface server running at http://52.91.86.137t:${port}`);
});