const express = require('express');
const path = require('path');
const axios = require('axios');
const cors = require('cors');

const app = express();
const port = 8080;

// Configuración de CORS
app.use(cors({
    origin: '*',  // Esto permite peticiones desde el puerto 8080
    methods: ['GET', 'POST', 'PUT', 'DELETE'],  // Puedes ajustar los métodos según sea necesario
    allowedHeaders: ['Content-Type', 'Authorization'],  // Cabeceras permitidas
}));


// Importing routes
const userRoutes = require('./routes/user-service');
const authRoutes = require('./routes/auth-service');
const sessionRoutes = require('./routes/session-service');
const catalogRoutes = require('./routes/catalog-service');

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

// Start server
app.listen(port, '0.0.0.0', () => {
    console.log(`Graphical interface server running at http://localhost:${port}`);
});