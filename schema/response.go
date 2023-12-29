package schema

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int64
}

type ErrorResponse struct {
	Message string `json:"message"`
}
