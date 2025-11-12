import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

export default function RegisterForm() {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [validationErrors, setValidationErrors] = useState({});

  const { register } = useAuth();
  const navigate = useNavigate();

  // Password validation rules
  const validatePassword = (password) => {
    const errors = {};
    
    if (password.length < 8) {
      errors.length = 'Password must be at least 8 characters';
    }
    
    if (!/[A-Z]/.test(password)) {
      errors.uppercase = 'Password must contain at least one uppercase letter';
    }
    
    if (!/[a-z]/.test(password)) {
      errors.lowercase = 'Password must contain at least one lowercase letter';
    }
    
    if (!/[0-9]/.test(password)) {
      errors.number = 'Password must contain at least one number';
    }
    
    if (!/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
      errors.special = 'Password must contain at least one special character (!@#$%^&*...)';
    }
    
    return errors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setValidationErrors({});

    // Validate password
    const passwordErrors = validatePassword(formData.password);
    if (Object.keys(passwordErrors).length > 0) {
      setValidationErrors(passwordErrors);
      return;
    }

    // Check password confirmation
    if (formData.password !== formData.confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    // Validate username
    if (formData.username.length < 3) {
      setError('Username must be at least 3 characters');
      return;
    }

    if (!/^[a-zA-Z0-9_]+$/.test(formData.username)) {
      setError('Username can only contain letters, numbers, and underscores');
      return;
    }

    setLoading(true);

    try {
      await register(formData.username, formData.email, formData.password);
      navigate('/dashboard');
    } catch (err) {
      setError(err.message || 'Registration failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });

    // Clear errors when user starts typing
    if (error) {
      setError('');
    }
    
    // Clear validation errors when user types
    if (name === 'password' && validationErrors) {
      setValidationErrors({});
    }
  };

  // Check password strength for visual feedback
  const getPasswordStrength = (password) => {
    const errors = validatePassword(password);
    const errorCount = Object.keys(errors).length;
    
    if (password.length === 0) return { strength: 0, label: '', color: '' };
    if (errorCount === 0) return { strength: 100, label: 'Strong', color: 'bg-green-500' };
    if (errorCount <= 2) return { strength: 60, label: 'Medium', color: 'bg-yellow-500' };
    return { strength: 30, label: 'Weak', color: 'bg-red-500' };
  };

  const passwordStrength = getPasswordStrength(formData.password);

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="bg-red-500/20 border border-red-500 text-red-200 px-4 py-3 rounded-lg animate-shake">
          <p className="text-sm font-medium">{error}</p>
        </div>
      )}

      <div>
        <label htmlFor="username" className="block text-sm font-medium mb-2 text-gray-200">
          Username
        </label>
        <input
          type="text"
          id="username"
          name="username"
          value={formData.username}
          onChange={handleChange}
          required
          minLength={3}
          maxLength={50}
          pattern="[a-zA-Z0-9_]+"
          className={`w-full px-4 py-2 bg-gray-700 border rounded-lg focus:outline-none focus:ring-2 transition-all text-white placeholder-gray-400 ${
            error 
              ? 'border-red-500 focus:ring-red-500 focus:border-red-500' 
              : 'border-gray-600 focus:ring-blue-500 focus:border-transparent'
          }`}
          placeholder="Choose a username"
          autoComplete="username"
        />
        <p className="text-xs text-gray-400 mt-1">
          3-50 characters, letters, numbers, and underscores only
        </p>
      </div>

      <div>
        <label htmlFor="email" className="block text-sm font-medium mb-2 text-gray-200">
          Email
        </label>
        <input
          type="email"
          id="email"
          name="email"
          value={formData.email}
          onChange={handleChange}
          required
          className={`w-full px-4 py-2 bg-gray-700 border rounded-lg focus:outline-none focus:ring-2 transition-all text-white placeholder-gray-400 ${
            error 
              ? 'border-red-500 focus:ring-red-500 focus:border-red-500' 
              : 'border-gray-600 focus:ring-blue-500 focus:border-transparent'
          }`}
          placeholder="Enter your email"
          autoComplete="email"
        />
      </div>

      <div>
        <label htmlFor="password" className="block text-sm font-medium mb-2 text-gray-200">
          Password
        </label>
        <input
          type="password"
          id="password"
          name="password"
          value={formData.password}
          onChange={handleChange}
          required
          minLength={8}
          className={`w-full px-4 py-2 bg-gray-700 border rounded-lg focus:outline-none focus:ring-2 transition-all text-white placeholder-gray-400 ${
            (error || Object.keys(validationErrors).length > 0)
              ? 'border-red-500 focus:ring-red-500 focus:border-red-500' 
              : 'border-gray-600 focus:ring-blue-500 focus:border-transparent'
          }`}
          placeholder="Create a password"
          autoComplete="new-password"
        />
        
        {/* Password strength indicator */}
        {formData.password && (
          <div className="mt-2">
            <div className="flex items-center justify-between mb-1">
              <span className="text-xs text-gray-400">Password strength:</span>
              <span className={`text-xs font-medium ${
                passwordStrength.strength === 100 ? 'text-green-400' :
                passwordStrength.strength >= 60 ? 'text-yellow-400' :
                'text-red-400'
              }`}>
                {passwordStrength.label}
              </span>
            </div>
            <div className="w-full bg-gray-600 rounded-full h-2">
              <div
                className={`h-2 rounded-full transition-all duration-300 ${passwordStrength.color}`}
                style={{ width: `${passwordStrength.strength}%` }}
              ></div>
            </div>
          </div>
        )}

        {/* Password requirements */}
        <div className="mt-2 space-y-1">
          <p className="text-xs text-gray-400">Password must contain:</p>
          <ul className="text-xs space-y-1">
            <li className={formData.password.length >= 8 ? 'text-green-400' : 'text-gray-400'}>
              ✓ At least 8 characters
            </li>
            <li className={/[A-Z]/.test(formData.password) ? 'text-green-400' : 'text-gray-400'}>
              ✓ One uppercase letter
            </li>
            <li className={/[a-z]/.test(formData.password) ? 'text-green-400' : 'text-gray-400'}>
              ✓ One lowercase letter
            </li>
            <li className={/[0-9]/.test(formData.password) ? 'text-green-400' : 'text-gray-400'}>
              ✓ One number
            </li>
            <li className={/[!@#$%^&*(),.?":{}|<>]/.test(formData.password) ? 'text-green-400' : 'text-gray-400'}>
              ✓ One special character (!@#$%^&*...)
            </li>
          </ul>
        </div>

        {/* Validation errors */}
        {Object.keys(validationErrors).length > 0 && (
          <div className="mt-2 bg-red-500/10 border border-red-500/50 rounded-lg p-2">
            {Object.values(validationErrors).map((error, index) => (
              <p key={index} className="text-xs text-red-400">• {error}</p>
            ))}
          </div>
        )}
      </div>

      <div>
        <label htmlFor="confirmPassword" className="block text-sm font-medium mb-2 text-gray-200">
          Confirm Password
        </label>
        <input
          type="password"
          id="confirmPassword"
          name="confirmPassword"
          value={formData.confirmPassword}
          onChange={handleChange}
          required
          minLength={8}
          className={`w-full px-4 py-2 bg-gray-700 border rounded-lg focus:outline-none focus:ring-2 transition-all text-white placeholder-gray-400 ${
            (error || (formData.confirmPassword && formData.password !== formData.confirmPassword))
              ? 'border-red-500 focus:ring-red-500 focus:border-red-500' 
              : 'border-gray-600 focus:ring-blue-500 focus:border-transparent'
          }`}
          placeholder="Confirm your password"
          autoComplete="new-password"
        />
        {formData.confirmPassword && formData.password !== formData.confirmPassword && (
          <p className="text-xs text-red-400 mt-1">Passwords do not match</p>
        )}
      </div>

      <button
        type="submit"
        disabled={loading}
        className="w-full btn-primary disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-105 active:scale-95"
      >
        {loading ? (
          <span className="flex items-center justify-center">
            <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Creating account...
          </span>
        ) : (
          'Create Account'
        )}
      </button>
    </form>
  );
}
