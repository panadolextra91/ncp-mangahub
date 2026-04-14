package models

// User represents the system identity, holding specific authorization roles and secure Bcrypt hashes.
type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         string
}
