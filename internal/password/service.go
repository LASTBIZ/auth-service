package password

import (
	"context"
	"lastbiz/auth-service/pkg/errors"
)

type Service struct {
	storage Storage
}

func NewPasswordService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s Service) CreatePassword(ctx context.Context, hash *Hash) error {
	err := s.storage.CreatePassword(*hash)
	return err
}

func (s Service) UpdatePassword(userID uint32, hash *Hash) error {
	if hash.Hash == "" {
		return errors.New("hash is empty")
	}
	err := s.storage.UpdatePassword(userID, hash.Hash)
	return err
}

func (s Service) DeletePassword(userID uint32) error {
	return s.storage.DeletePassword(userID)
}

func (s Service) GetHash(userID uint32) (*Hash, error) {
	return s.storage.GetHash(userID)
}
