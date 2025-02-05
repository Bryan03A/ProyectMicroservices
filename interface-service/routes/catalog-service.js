const express = require('express');
const axios = require('axios');
const router = express.Router();

// URL of the catalog-service microservice
const CATALOG_SERVICE_URL = process.env.CATALOG_SERVICE_URL || 'http://54.214.66.77:5003';

// Helper function to handle requests to the microservice
const fetchFromCatalogService = async (method, url, data = null, token = null) => {
    try {
        const config = { method, url, data, headers: {} };
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        const response = await axios(config);
        return { success: true, data: response.data };
    } catch (error) {
        console.error(`Error in catalog-service request: ${error.message}`);
        return { success: false, error: error.response?.data || 'Unknown error' };
    }
};

// Route to get all models
router.get('/models', async (req, res) => {
    const { success, data, error } = await fetchFromCatalogService('get', `${CATALOG_SERVICE_URL}/models`);
    if (success) {
        res.status(200).json(data);
    } else {
        res.status(500).json({ message: 'Error fetching models', error });
    }
});

// Route to get models by user ID
router.get('/models/user/:userId', async (req, res) => {
    const userId = req.params.userId;
    const { success, data, error } = await fetchFromCatalogService('get', `${CATALOG_SERVICE_URL}/models/user/${userId}`);
    if (success) {
        res.status(200).json(data);
    } else {
        res.status(500).json({ message: 'Error fetching user models', error });
    }
});

// Route to add a new model
router.post('/models', async (req, res) => {
    const token = req.headers.authorization?.split(" ")[1];  // Extract token from headers
    const { success, data, error } = await fetchFromCatalogService('post', `${CATALOG_SERVICE_URL}/models`, req.body, token);
    if (success) {
        res.status(201).json(data);
    } else {
        res.status(500).json({ message: 'Error adding model', error });
    }
});

// Route to get a specific model by name
router.get('/models/:modelName', async (req, res) => {
    const modelName = req.params.modelName;
    const { success, data, error } = await fetchFromCatalogService('get', `${CATALOG_SERVICE_URL}/models/${modelName}`);
    if (success) {
        res.status(200).json(data);
    } else {
        res.status(404).json({ message: 'Model not found', error });
    }
});

// Route to delete a model
router.delete('/models/:modelName', async (req, res) => {
    const modelName = req.params.modelName;
    const token = req.headers.authorization?.split(" ")[1];  // Extract token from headers
    const { success, data, error } = await fetchFromCatalogService('delete', `${CATALOG_SERVICE_URL}/models/${modelName}`, null, token);
    if (success) {
        res.status(200).json(data);
    } else {
        res.status(404).json({ message: 'Error deleting model', error });
    }
});

module.exports = router;