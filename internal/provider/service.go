package provider

import (
	"context"
	"lastbiz/auth-service/internal/auth"
)

type Service struct {
	storage Storage
}

func NewProviderService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s Service) CreateProvider(ctx context.Context, provider *auth.OAuthProvider) error {
	providerStorageMap, err := provider.ToMap()
	if err != nil {
		return err
	}
	err = s.storage.CreateProvider(ctx, providerStorageMap)
	return err
}

func (s Service) UpdateProvider(ctx context.Context, userID uint32, provider string, providerStruct *auth.OAuthProvider) error {
	providerStorageMap, err := providerStruct.ToMap()
	if err != nil {
		return err
	}
	err = s.storage.UpdateProvider(ctx, provider, userID, providerStorageMap)
	return err
}

func (s Service) DeletePassword(ctx context.Context, userID uint32, provider string) error {
	return s.storage.DeleteProvider(ctx, userID, provider)
}

func (s Service) GetProvider(ctx context.Context, userID uint32, provider string) (*auth.OAuthProvider, error) {
	return s.storage.GetProvider(ctx, userID, provider)
}
