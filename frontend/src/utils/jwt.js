/**
 * Decode JWT token payload
 * @param {string} token - JWT token string
 * @returns {Object|null} Decoded payload or null if invalid
 */
export function decodeJWT(token) {
  if (!token) return null;
  
  try {
    const parts = token.split('.');
    if (parts.length !== 3) return null;
    
    // Decode the payload (second part)
    const payload = parts[1];
    const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'));
    return JSON.parse(decoded);
  } catch (error) {
    console.error('Failed to decode JWT:', error);
    return null;
  }
}

/**
 * Extract user info from JWT token
 * @param {string} token - JWT token string
 * @returns {Object|null} User info with user_id and username
 */
export function getUserFromToken(token) {
  const payload = decodeJWT(token);
  if (!payload) return null;
  
  return {
    id: payload.user_id,
    username: payload.username,
  };
}

/**
 * Check if JWT token is expired
 * @param {string} token - JWT token string
 * @returns {boolean} True if expired
 */
export function isTokenExpired(token) {
  const payload = decodeJWT(token);
  if (!payload || !payload.exp) return true;
  
  // exp is in seconds, Date.now() is in milliseconds
  return payload.exp * 1000 < Date.now();
}
