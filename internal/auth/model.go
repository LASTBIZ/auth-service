package auth

import "time"

type PasswordHash struct {
	ID     int
	UserID uint32
	Hash   string
}

type OAuthProvider struct {
	ID            int
	UserID        uint32
	OAuthProvider string
	AccessToken   string
	RefreshToken  string
	ExpiryDate    time.Time
}
