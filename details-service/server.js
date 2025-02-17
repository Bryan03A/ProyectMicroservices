const express = require('express');
const path = require('path');
const axios = require('axios');
const cors = require('cors');

const app = express();
const port = 5015;

// CORS Configuration
app.use(cors({
    origin: '*',  // This allows requests from all origins
    methods: ['GET', 'POST', 'PUT', 'DELETE'],  // You can adjust the methods as needed
    allowedHeaders: ['Content-Type', 'Authorization'],  // Allowed headers
}));

// Importing routes
const userRoutes = require('./routes/user-service');
const authRoutes = require('./routes/auth-service');
const sessionRoutes = require('./routes/session-service');
const catalogRoutes = require('./routes/catalog-service');
const searchService = require("./routes/search-service");

// Serving static files (like index.html)
app.use(express.static(path.join(__dirname, 'public')));

// Middleware for parsing JSON requests
app.use(express.json());

// Use the user routes
app.use(userRoutes);

// Use the authentication routes
app.use(authRoutes);

// Use session routes
app.use(sessionRoutes);

// Use the catalog routes
app.use(catalogRoutes);

// Route to handle search requests
app.use('/search', require('./routes/search-service'));

// Route to serve the chat page
app.get('/chat', (req, res) => {
    res.sendFile(path.join(__dirname, 'public', 'chat.html'));
});

// Start the server
app.listen(port, '0.0.0.0', () => {
    console.log(`Graphical interface server running at http://localhost:${port}`);
});