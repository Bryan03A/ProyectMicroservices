const express = require("express");
const redis = require("redis");
const router = express.Router();

// Configurar Redis
const redisClient = redis.createClient({
    socket: {
        host: "localhost",  // Cambia esto si usas Docker (puede ser "redis-server")
        port: 6379
    }
});

// Conectar a Redis con manejo de errores
(async () => {
    try {
        await redisClient.connect();
        console.log("✅ Conectado a Redis");
    } catch (err) {
        console.error("❌ Error al conectar a Redis:", err);
    }
})();

// Middleware para manejar errores de Redis
redisClient.on("error", (err) => console.error("❌ Redis Error:", err));

// Guardar búsqueda en Redis para un usuario específico
router.post("/save-search", async (req, res) => {
    const { query, creator, username, firstModelName } = req.body;  // Ahora agregamos firstModelName

    if (!query || !username) {
        return res.status(400).send("Falta la consulta de búsqueda o el nombre de usuario");
    }

    try {
        // Guardamos la búsqueda junto con el nombre del primer modelo en una lista de historial para el usuario
        const userSearchKey = `search-history:${username}`;
        const searchData = {
            query,
            creator,
            firstModelName,  // Agregamos el nombre del primer modelo
            timestamp: Date.now()
        };
        await redisClient.lPush(userSearchKey, JSON.stringify(searchData));

        res.status(200).send("Búsqueda guardada en Redis");
    } catch (error) {
        console.error("Error al guardar búsqueda en Redis:", error);
        res.status(500).send("Error al guardar búsqueda");
    }
});

// Obtener las últimas búsquedas guardadas en Redis para un usuario específico
router.get("/recent-searches", async (req, res) => {
    const { username } = req.query;  // Obtenemos el username de la consulta

    if (!username) {
        return res.status(400).send("Falta el nombre de usuario");
    }

    try {
        const userSearchKey = `search-history:${username}`;
        const searches = await redisClient.lRange(userSearchKey, 0, 9);
        const parsedSearches = searches.map((search) => JSON.parse(search));

        res.json(parsedSearches);
    } catch (error) {
        console.error("Error al obtener búsquedas recientes:", error);
        res.status(500).send("Error al obtener búsquedas recientes");
    }
});

module.exports = router;