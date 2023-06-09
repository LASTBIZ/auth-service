package data

import "time"

type RefreshSession struct {
	ID           uint64
	userID       uint64
	refreshToken string
	ua           string
	ip           string
	expiresIn    time.Time
	createdAt    time.Time
}
