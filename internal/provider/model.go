package provider

import (
	"time"
)

type OAuthProvider struct {
	ID           uint32 `mapper:"id"`
	UserID       uint32 `mapper:"user_id"`
	Provider     string `mapper:"provider"`
	AccessToken  string `mapper:"access_token"`
	RefreshToken string `mapper:"refresh_token"`
	TokenType    string
	ExpiryDate   time.Time `mapper:"expiry_date"`
}
