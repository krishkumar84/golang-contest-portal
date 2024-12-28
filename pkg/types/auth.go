package types

type LoginRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required"`
}

type TokenResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
}