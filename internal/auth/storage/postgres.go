package storage

import (
	"gorm.io/gorm"
	"lastbiz/auth-service/internal/auth"
)

type storage struct {
	db gorm.DB
}

func NewAuthStorage(db gorm.DB) auth.Storage {
	return &storage{
		db: db,
	}
}

func (s storage) CreatePassword(hash auth.PasswordHash) error {
	err := s.db.Create(&hash).Error
	return err
}

func (s storage) UpdatePassword(hash auth.PasswordHash) error {
	err := s.db.
		Model(&auth.PasswordHash{}).
		Where(&auth.PasswordHash{UserID: hash.UserID}).
		Update("hash", hash.Hash).
		Error
	return err
}

func (s storage) CreateProvider(provider auth.OAuthProvider) error {
	err := s.db.Create(&provider).Error
	return err
}

func (s storage) UpdateProvider(provider auth.OAuthProvider) error {
	err := s.db.
		Save(provider).
		Error
	return err
}

func (s storage) DeletePassword(userID uint32) error {
	err := s.db.Model(&auth.PasswordHash{}).
		Where(&auth.PasswordHash{UserID: userID}).
		Delete(&auth.PasswordHash{}).Error
	return err
}

func (s storage) GetHash(userID uint32) (auth.PasswordHash, error) {
	var hash auth.PasswordHash
	err := s.db.Model(&auth.PasswordHash{}).
		Where(&auth.PasswordHash{UserID: userID}).
		First(&hash).
		Error
	return hash, err
}

func (s storage) DeleteProvider(userID uint32, provider string) error {
	err := s.db.Model(&auth.OAuthProvider{}).
		Where(&auth.OAuthProvider{UserID: userID, OAuthProvider: provider}).
		Delete(&auth.OAuthProvider{}).
		Error
	return err
}

func (s storage) GetProvider(userID uint32, provider string) (auth.OAuthProvider, error) {
	var _provider auth.OAuthProvider
	err := s.db.Model(&auth.OAuthProvider{}).
		Where(&auth.OAuthProvider{UserID: userID, OAuthProvider: provider}).
		First(&_provider).
		Error
	return _provider, err
}
