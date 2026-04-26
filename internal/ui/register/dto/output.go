package dto

type RegisterOutput struct {
	AccessToken           string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	RefreshToken          string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.refresh"`
	AccessTokenExpiresAt  string `json:"access_token_expires_at" example:"2026-04-27T12:00:00Z"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at" example:"2026-05-27T12:00:00Z"`
} //	@name	RegisterOutput
