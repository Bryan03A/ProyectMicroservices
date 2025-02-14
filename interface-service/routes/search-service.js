const express = require('express');
const axios = require('axios');
const router = express.Router();

// Ruta para búsqueda avanzada
router.get('/', async (req, res) => {
    const query = req.query.q ? req.query.q.toLowerCase() : '';  // Convertir a minúsculas
    const creator = req.query.creator ? req.query.creator.toLowerCase() : ''; // Convertir creator a minúsculas
    const sortBy = req.query.sort === 'price' ? 'price' : 'name.keyword'; // Asegurar orden por un campo keyword
    const sortOrder = req.query.order === 'desc' ? 'desc' : 'asc'; // Dirección de orden

    try {
        const mustQueries = [];

        if (query) {
            mustQueries.push({
                bool: {
                    should: [
                        { match_phrase_prefix: { name: query } }, // Buscar coincidencias al inicio en name
                        { match: { description: { query, fuzziness: "AUTO" } } } // Búsqueda flexible en description
                    ]
                }
            });
        }

        if (creator) {
            mustQueries.push({
                match_phrase_prefix: { created_by: creator } // Coincidencias parciales en el creador
            });
        }

        const esQuery = {
            query: {
                bool: {
                    must: mustQueries
                }
            },
            sort: [
                { [sortBy]: { order: sortOrder } } // Ordenar por un campo que permita sorting
            ]
        };

        // Realizar la consulta a Elasticsearch
        const response = await axios.post('http://localhost:9200/models/_search', esQuery, {
            headers: { 'Content-Type': 'application/json' }
        });

        res.json(response.data.hits.hits);
    } catch (error) {
        console.error('Error en la búsqueda:', error.response?.data || error.message);
        res.status(500).send('Error al realizar la búsqueda');
    }
});

module.exports = router;