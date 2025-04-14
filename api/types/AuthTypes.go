package types

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token string
}

type RegisterRequest struct {
	Username string
	Password string
	Email    string
}

type LogoutRequest struct {
	Token string
}
