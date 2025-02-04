// Configuración de URLs para la API
const config = {
    apiUrl: window.location.hostname === 'localhost' ? 'http://localhost:5001' : 'http://auth-service:5001',
    sessionApiUrl: window.location.hostname === 'localhost' ? 'http://localhost:5004' : 'http://session-service:5004',
    catalogApiUrl: window.location.hostname === 'localhost' ? 'http://localhost:5003' : 'http://catalog-service:5003'
};

// Hacer disponible la configuración en el ámbito global
window.config = config;
