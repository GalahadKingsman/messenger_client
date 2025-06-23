package models

import "net/http"

type CreateUserRequest struct {
	Login     string `json:"login"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type CreateUserResponse struct {
	Success string `json:"success"`
	ID      int64  `json:"-"`
}

type Client struct {
	APIGatewayURL string
	HTTPClient    *http.Client
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	UserID  int64  `json:"user_id,omitempty"`
}
