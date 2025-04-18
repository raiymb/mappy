package dto

// RegisterDTO is the payload for POST /auth/register
type RegisterDTO struct {
	Email    string `json:"email"    validate:"required,email"`
	Name     string `json:"name"     validate:"required,min=2"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginDTO is the payload for POST /auth/login
type LoginDTO struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
