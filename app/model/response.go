package model

// ErrorResponse represents a standard error payload
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"invalid request"`
}

// TokenResponse represents the token response
type TokenResponse struct {
	AccessToken string `json:"access_token" example:"<jwt-token>"`
	Msg         string `json:"msg" example:"Token will be expired within 15 minutes"`
}
