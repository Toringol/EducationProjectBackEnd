package models

// User - user model in DB
type User struct {
	ID       int64  `json:"id,omitempty"`
	Email    string `json:"email"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	Avatar   string `json:"avatar"`
	Role     int    `json:"role"`
}
