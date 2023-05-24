package data

import (
	"auth-service/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type Hash struct {
	ID     uint32
	UserID uint32
	Hash   string
}

type hashRepo struct {
	data *Data
	log  *log.Helper
}

func NewHashRepo(data *Data, logger log.Logger) biz.HashRepo {
	return &hashRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (h hashRepo) CreateHash(ctx context.Context, hb *biz.Hash) (*biz.Hash, error) {
	var hash Hash
	result := h.data.db.Where(&Hash{UserID: hb.UserID}).First(&hash)
	if result.RowsAffected == 1 {
		return nil, errors.New(500, "USER_EXISTS", "user is exists")
	}

	hash.UserID = hb.UserID
	hash.Hash = hb.Hash

	res := h.data.db.Create(&hash)
	if res.Error != nil {
		return nil, errors.New(500, "CREATE_USER_ERROR", "error create user")
	}

	hashInfoRes := modelToResponse(hash)
	return hashInfoRes, nil
}

func (h hashRepo) UpdateHash(ctx context.Context, hash *biz.Hash) (bool, error) {
	var hashInfo Hash
	result := h.data.db.Where(&Hash{UserID: hash.UserID}).First(&hashInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, errors.NotFound("USER_NOT_FOUND", "user not found")
	}

	if result.RowsAffected == 0 {
		return false, errors.NotFound("USER_NOT_FOUND", "rows null")
	}

	hashInfo.Hash = hash.Hash

	if err := h.data.db.Save(&hashInfo).Error; err != nil {
		return false, errors.New(500, "USER_UPDATE_ERROR", "update user error")
	}

	return true, nil
}

func (h hashRepo) DeleteHashByUserId(ctx context.Context, userId uint32) (bool, error) {
	var hashInfo Hash
	result := h.data.db.Where(&Hash{UserID: userId}).First(&hashInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, errors.NotFound("USER_NOT_FOUND", "user not found")
	}

	if result.RowsAffected == 0 {
		return false, errors.NotFound("USER_NOT_FOUND", "rows null")
	}

	if err := h.data.db.Delete(&hashInfo).Error; err != nil {
		return false, errors.New(500, "USER_DELETE_ERROR", "delete user error")
	}

	return true, nil
}

func (h hashRepo) GetHashByUserId(ctx context.Context, userId uint32) (*biz.Hash, error) {
	var hashInfo Hash
	if err := h.data.db.Where(&Hash{UserID: userId}).First(&hashInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
		}

		return nil, errors.New(500, "USER_NOT_FOUND", err.Error())
	}

	re := modelToResponse(hashInfo)
	return re, nil
}

func modelToResponse(hash Hash) *biz.Hash {
	hashInfoRsp := &biz.Hash{
		ID:     hash.ID,
		UserID: hash.UserID,
		Hash:   hash.Hash,
	}
	return hashInfoRsp
}
