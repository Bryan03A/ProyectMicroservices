// Configuration object for API URLs based on the environment
const config = {
    // apiUrl: The URL for the main API, checking if the current environment is 'localhost' or production.
    // If it's 'localhost', it uses the local development URL (localhost); otherwise, it uses the production URL.
    apiUrl: window.location.hostname === 'localhost' ? 'http://52.91.86.137:5001' : 'http://52.91.86.137:5001',

    // sessionApiUrl: The URL for the session-related API, again checking for 'localhost' or production.
    // If it's 'localhost', it points to a different port for local development; otherwise, it uses the production server.
    sessionApiUrl: window.location.hostname === 'localhost' ? 'http://52.91.86.137:5004' : 'http://52.91.86.137:5004',

    // catalogApiUrl: Similar to the other URLs, this is for accessing the catalog API, adjusting based on the environment.
    catalogApiUrl: window.location.hostname === 'localhost' ? 'http://52.91.86.137:5003' : 'http://52.91.86.137:5003'
};

// Expose the config object globally so it can be accessed anywhere in the client-side code
window.config = config;