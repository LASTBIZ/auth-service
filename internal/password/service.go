package password

import (
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

func (s Service) CreatePassword(hash *Hash) error {
	err := s.storage.CreatePassword(*hash)
	return err
}

func (s Service) UpdatePassword(userID uint32, hash *Hash) error {
	if hash.Hash == "" {
		return errors.New("hash not found")
	}
	err := s.storage.UpdatePassword(userID, hash.Hash)
	return err
}

func (s Service) DeletePassword(userID uint32) error {
	return s.storage.DeletePassword(userID)
}

func (s Service) GetHash(userID uint32) (*Hash, error) {
	if userID == 0 {
		return nil, errors.New("hash not found")
	}
	return s.storage.GetHash(userID)
}
