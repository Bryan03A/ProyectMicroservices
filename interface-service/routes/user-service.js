const express = require('express');
const axios = require('axios');

const router = express.Router();
<<<<<<< HEAD
const USER_SERVICE_URL = process.env.USER_SERVICE_URL || 'http://98.84.113.50:5002';
=======
const USER_SERVICE_URL = process.env.USER_SERVICE_URL || 'http://localhost:5002';
>>>>>>> b6eb336 (test)

// Route to handle new user registration
router.post('/register', async (req, res) => {
    try {
        const response = await axios.post(`${USER_SERVICE_URL}/register`, req.body);
        res.status(201).json(response.data);
    } catch (error) {
        console.error(error.response?.data || error.message);
        res.status(500).json({ message: 'Error registering user' });
    }
});

module.exports = router;