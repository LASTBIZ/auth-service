package provider

import (
	"gorm.io/gorm"
	"time"
)

type Storage struct {
	db gorm.DB
}

func NewProviderStorage(db gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s Storage) CreateProvider(_provider OAuthProvider) error {
	return s.db.Create(&_provider).Error
}

func (s Storage) UpdateProvider(
	provider string,
	userID uint32,
	accessToken,
	refreshToken string,
	expiryDate time.Time,
) error {
	return s.db.Model(&OAuthProvider{}).
		Where("user_id = ? and provider = ?", userID, provider).
		Updates(map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expiry_date":   expiryDate,
		}).Error
}

func (s Storage) DeleteProvider(userID uint32, provider string) error {
	return s.db.
		Where("user_id = ? and provider = ?", userID, provider).
		Delete(&OAuthProvider{}).
		Error
}

func (s Storage) GetProvider(userID uint32, provider string) (*OAuthProvider, error) {
	var getProvider OAuthProvider
	result := s.db.
		Where("user_id = ? and provider = ?", userID, provider).
		First(&getProvider)
	return &getProvider, result.Error
}

func (s Storage) CheckProvider(userID uint32) (*OAuthProvider, error) {
	var getProvider OAuthProvider
	result := s.db.
		Where("user_id = ?", userID).
		First(&getProvider)
	return &getProvider, result.Error
}
