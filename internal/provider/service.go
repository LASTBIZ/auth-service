package provider

import (
	"context"
)

type Service struct {
	storage Storage
}

func NewProviderService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s Service) CreateProvider(ctx context.Context, provider *OAuthProvider) error {
	providerStorageMap, err := provider.ToMap()
	if err != nil {
		return err
	}
	err = s.storage.CreateProvider(ctx, providerStorageMap)
	return err
}

func (s Service) UpdateProvider(ctx context.Context, userID uint32, provider string, providerStruct *OAuthProvider) error {
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

func (s Service) GetProvider(ctx context.Context, userID uint32, provider string) (*OAuthProvider, error) {
	return s.storage.GetProvider(ctx, userID, provider)
}
