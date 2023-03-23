// API base URL
export const API_BASE_URL = process.env.REACT_APP_API_ROOT || '/api'

// OIDC authority URL
export const AUTHORITY = process.env.REACT_APP_AUTHORITY || 'none'

// OIDC client ID
export const CLIENT_ID = process.env.REACT_APP_CLIENT_ID || 'readflow.app'

// Unauthenticated user redirect
export const REDIRECT_URL = process.env.REACT_APP_REDIRECT_URL || '/login'

// VERSION
export const VERSION = process.env.REACT_APP_VERSION || 'snapshot'
