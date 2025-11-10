package auth

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents an authentication response with token and user info
type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"` // Will be *database.User
}
