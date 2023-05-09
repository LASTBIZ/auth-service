package auth

import "time"

type Service struct {
}

func NewAuthService() Service {
	return Service{}
}

func (s Service) CreatePassword(hash string, userID uint32) error {

}

func (s Service) UpdatePassword(hash string, userID uint32) error {

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
