package models

// User represents the system identity, holding specific authorization roles and secure Bcrypt hashes.
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
}
