// Configuración de URLs para la API
const config = {
    apiUrl: window.location.hostname === 'localhost' ? 'http://54.214.66.77:5001' : 'http://54.214.66.77:5001',
    sessionApiUrl: window.location.hostname === 'localhost' ? 'http://54.214.66.77:5004' : 'http://54.214.66.77:5004',
    catalogApiUrl: window.location.hostname === 'localhost' ? 'http://54.214.66.77:5003' : 'http://54.214.66.77:5003'
};

// Hacer disponible la configuración en el ámbito global
window.config = config;
