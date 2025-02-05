const express = require('express');
const axios = require('axios');
const router = express.Router();

const SESSION_SERVICE_URL = process.env.SESSION_SERVICE_URL || 'http://54.214.66.77:5004';  // URL of your session-service

// Route to create a session
router.post('/create-session', async (req, res) => {
    try {
        const token = req.headers.authorization?.split(' ')[1];
        if (!token) {
            return res.status(403).json({ message: 'No token provided' });
        }

        // Send token to session-service to create a session
        const sessionResponse = await axios.post(`${SESSION_SERVICE_URL}/api/sessions`, {
            data: {}  // Session data, you can add preferences or other data here
        }, {
            headers: { Authorization: `Bearer ${token}` }
        });

        res.status(201).json(sessionResponse.data);
    } catch (error) {
        console.error(error.response?.data || error.message);
        res.status(500).json({ message: 'Error creating session', error: error.message });
    }
});

module.exports = router;