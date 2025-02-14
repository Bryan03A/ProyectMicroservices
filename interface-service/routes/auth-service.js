const express = require('express');
const axios = require('axios');

const router = express.Router();
const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:5001';

// Route to handle user login
router.post('/login', async (req, res) => {
    try {
        const response = await axios.post(`${AUTH_SERVICE_URL}/login`, req.body);
        res.json(response.data);
    } catch (error) {
        console.error(error.response?.data || error.message);
        res.status(401).json({ message: 'Invalid credentials' });
    }
});

// Route to obtain user profile data
router.get('/profile', async (req, res) => {
    try {
        const token = req.headers.authorization?.split(' ')[1];
        if (!token) {
            return res.status(403).json({ message: 'No token provided' });
        }

        const response = await axios.get(`${AUTH_SERVICE_URL}/profile`, {
            headers: { Authorization: `Bearer ${token}` }
        });

        res.json(response.data);
    } catch (error) {
        console.error(error.response?.data || error.message);
        res.status(500).json({ message: 'Error fetching profile' });
    }
});

module.exports = router;
