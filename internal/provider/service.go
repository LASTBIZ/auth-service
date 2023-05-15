package provider

import (
	"lastbiz/auth-service/internal/utils"
	"lastbiz/auth-service/pkg/errors"
	"time"
)

type Service struct {
	storage        Storage
	providers      map[string]Provider
	providersNames []string
}

func NewProviderService(storage Storage, providers map[string]Provider) *Service {
	keys := make([]string, 0, len(providers))
	for k := range providers {
		keys = append(keys, k)
	}
	return &Service{
		storage:        storage,
		providers:      providers,
		providersNames: keys,
	}
}

func (s Service) GetProviderByName(name string) (Provider, error) {
	val, ok := s.providers[name]
	if !ok {
		return nil, errors.New("provider not found")
	}
	return val, nil
}

func (s Service) CreateProvider(provider *OAuthProvider) error {
	if utils.Contains(s.providersNames, provider.Provider) {
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
	if utils.Contains(s.providersNames, provider) {
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
