package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenPayload struct {
	ID     uint
	Email  string
	IsRoot bool
}

type RefreshTokenRequest struct {
	AccessToken string `json:"access_token"`
}

type RefreshTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
}
