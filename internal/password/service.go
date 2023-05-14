package password

import (
	"context"
	"lastbiz/auth-service/internal/auth"
)

type Service struct {
	storage Storage
}

func NewPasswordService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s Service) CreatePassword(ctx context.Context, hash *auth.PasswordHash) error {
	hashStorageMap, err := hash.ToMap()
	if err != nil {
		return err
	}
	err = s.storage.CreatePassword(ctx, hashStorageMap)
	return err
}

func (s Service) UpdatePassword(ctx context.Context, id uint32, hash *auth.PasswordHash) error {
	hashStorageMap, err := hash.ToMap()
	if err != nil {
		return err
	}
	err = s.storage.UpdatePassword(ctx, id, hashStorageMap)
	return err
}

func (s Service) DeletePassword(ctx context.Context, id uint32) error {
	return s.storage.DeletePassword(ctx, id)
}

func (s Service) GetHash(ctx context.Context, userID uint32) (*auth.PasswordHash, error) {
	return s.storage.GetHash(ctx, userID)
}
