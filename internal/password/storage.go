package password

import (
	"gorm.io/gorm"
)

type Storage struct {
	db gorm.DB
}

func NewPasswordStorage(db gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s Storage) CreatePassword(hashCreate Hash) error {
	return s.db.Create(&hashCreate).Error
}

func (s Storage) UpdatePassword(userID uint32, hash string) error {
	return s.db.Model(&Hash{}).Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"hash": hash,
		}).Error
}

func (s Storage) DeletePassword(userID uint32) error {
	return s.db.Where("user_id = ?", userID).Delete(&Hash{}).Error
}

func (s Storage) GetHash(userID uint32) (*Hash, error) {
	var getHash Hash
	result := s.db.Where("user_id = ?", userID).
		First(&getHash)
	return &getHash, result.Error
}
