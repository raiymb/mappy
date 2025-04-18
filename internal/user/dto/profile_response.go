package dto

// ProfileResponse returned by GET /users/me
type ProfileResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}
