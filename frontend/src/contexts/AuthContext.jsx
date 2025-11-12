import { createContext, useState, useEffect, useContext } from 'react';
import authService from '../services/auth.service';
import { getUserFromToken } from '../utils/jwt';

export const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Initialize auth state from localStorage
  useEffect(() => {
    const initAuth = async () => {
      const token = authService.getStoredToken();
      
      if (token) {
        try {
          // First, get user info from token (includes username)
          const tokenUser = getUserFromToken(token);
          
          // Then fetch full user data from API
          const apiResponse = await authService.getCurrentUser();
          const currentUser = apiResponse.user || apiResponse;
          
          // Merge token data with API data to ensure username is present
          setUser({
            ...currentUser,
            username: tokenUser?.username || currentUser.username,
          });
        } catch (err) {
          // Token is invalid, clear storage
          authService.logout();
          setUser(null);
        }
      }
      
      setLoading(false);
    };

    initAuth();
  }, []);

  /**
   * Register a new user
   */
  const register = async (username, email, password) => {
    try {
      setError(null);
      const data = await authService.register(username, email, password);
      setUser(data.user);
      return data;
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  /**
   * Login with username and password
   */
  const login = async (username, password) => {
    try {
      setError(null);
      const data = await authService.login(username, password);
      setUser(data.user);
      return data;
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  /**
   * Logout current user
   */
  const logout = () => {
    authService.logout();
    setUser(null);
    setError(null);
  };

  /**
   * Update user data (e.g., after earning coins)
   */
  const updateUser = (updates) => {
    setUser(prev => ({ ...prev, ...updates }));
    localStorage.setItem('user', JSON.stringify({ ...user, ...updates }));
  };

  const value = {
    user,
    loading,
    error,
    register,
    login,
    logout,
    updateUser,
    isAuthenticated: !!user,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

/**
 * Custom hook to use auth context
 */
export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
