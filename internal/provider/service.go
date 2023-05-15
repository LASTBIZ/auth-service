package provider

import (
	"lastbiz/auth-service/internal/utils"
	"lastbiz/auth-service/pkg/errors"
	"time"
)

type Service struct {
	storage Storage
}

var providers []string

func NewProviderService(storage Storage) *Service {
	providers = []string{"google", "facebook"}
	return &Service{
		storage: storage,
	}
}

func (s Service) CreateProvider(provider *OAuthProvider) error {
	if utils.Contains(providers, provider.Provider) {
		return errors.New("provider not found")
	}
	err := s.storage.CreateProvider(*provider)
	return err
}

func (s Service) UpdateProvider(
	provider string,
	userID uint32,
	accessToken string,
	refreshToken string,
	expiryDate time.Time,
) error {
	if utils.Contains(providers, provider) {
		return errors.New("provider not found")
	}
	err := s.storage.UpdateProvider(
		provider,
		userID,
		accessToken,
		refreshToken,
		expiryDate)
	return err
}

func (s Service) DeletePassword(userID uint32, provider string) error {
	return s.storage.DeleteProvider(userID, provider)
}

func (s Service) GetProvider(userID uint32, provider string) (*OAuthProvider, error) {
	return s.storage.GetProvider(userID, provider)
}
