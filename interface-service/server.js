const express = require('express');
const path = require('path');
const axios = require('axios');
const cors = require('cors');

const app = express();
const port = 8080;

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
<<<<<<< HEAD
=======
const searchService = require("./routes/search-service");
const redisService = require("./routes/redis-service");
>>>>>>> b6eb336 (test)

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

<<<<<<< HEAD
// Start server
app.listen(port, '0.0.0.0', () => {
    console.log(`Graphical interface server running at http://3.82.92.84:${port}`);
=======
// Ruta para manejar las búsquedas
app.use('/search', require('./routes/search-service'));

app.use("/search", searchService);
app.use("/redis", redisService);

// Start server
app.listen(port, '0.0.0.0', () => {
    console.log(`Graphical interface server running at http://localhost:${port}`);
>>>>>>> b6eb336 (test)
});