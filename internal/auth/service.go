package auth

import (
	"time"
)

type Service struct {
	storage Storage
}

func NewAuthService() Service {
	return Service{}
}

func (s Service) CreatePassword(hash string, userID uint32) error {
	var _hash PasswordHash
	_hash = PasswordHash{UserID: userID, Hash: hash}
	return s.storage.CreatePassword(_hash)
}

func (s Service) UpdatePassword(hash string, userID uint32) error {

	var _hash PasswordHash
	_hash = PasswordHash{UserID: userID, Hash: hash}
	return s.storage.UpdatePassword(_hash)
}

func (s Service) GetHash(userID uint32) (string, error) {

}

func (s Service) CreateProvider(
	provider string,
	accessToken string,
	refreshToken string,
	expiryDate time.Time,
	userID uint32) error {

}

func (s Service) UpdateProvider(
	provider string,
	accessToken string,
	refreshToken string,
	expiryDate time.Time,
	userID uint32) error {

}
