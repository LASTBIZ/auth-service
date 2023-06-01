package data

import (
	"auth-service/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type Provider struct {
	ID           uint32
	UserID       uint32
	Email        string
	Provider     string
	AccessToken  string
	RefreshToken string
	TokenType    string
}

type providerRepo struct {
	data *Data
	log  *log.Helper
}

func NewProviderRepo(data *Data, logger log.Logger) biz.ProviderRepo {
	return &providerRepo{data: data, log: log.NewHelper(logger)}
}

func (p providerRepo) CreateProvider(ctx context.Context, pr *biz.Provider) (*biz.Provider, error) {
	var provider Provider
	result := p.data.db.Where(&Provider{UserID: pr.UserID}).Or(&Provider{Email: pr.Email}).First(&provider)
	if result.RowsAffected == 1 {
		return nil, errors.InternalServer("USER_EXISTS", "user is exists")
	}

	provider.Email = pr.Email
	provider.UserID = pr.UserID
	provider.Provider = pr.Provider
	provider.AccessToken = pr.AccessToken
	provider.RefreshToken = pr.RefreshToken
	provider.TokenType = pr.TokenType

	res := p.data.db.Create(&provider)
	if res.Error != nil {
		return nil, errors.InternalServer("CREATE_USER_ERROR", "error create user")
	}

	providerInfoRes := p.modelToResponse(provider)
	return providerInfoRes, nil
}

func (p providerRepo) CreateState(ctx context.Context) (string, error) {

}

func (p providerRepo) CheckState(ctx context.Context, state string) (string, error) {

}

func (p providerRepo) UpdateProvider(ctx context.Context, pr *biz.Provider) (bool, error) {
	var provider Provider
	result := p.data.db.Where(&Provider{UserID: pr.UserID}).Or(&Provider{Email: pr.Email}).First(&provider)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, errors.NotFound("USER_NOT_FOUND", "user not found")
	}

	if result.RowsAffected == 0 {
		return false, errors.NotFound("USER_NOT_FOUND", "rows null")
	}

	provider.AccessToken = pr.AccessToken
	provider.RefreshToken = pr.RefreshToken

	if err := p.data.db.Save(&provider).Error; err != nil {
		return false, errors.New(500, "USER_UPDATE_ERROR", "update user error")
	}

	return true, nil
}

func (p providerRepo) DeleteProviderByUserId(ctx context.Context, userId uint32) (bool, error) {
	var provider Provider
	result := p.data.db.Where(&Provider{UserID: userId}).First(&provider)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, errors.NotFound("USER_NOT_FOUND", "user not found")
	}

	if result.RowsAffected == 0 {
		return false, errors.NotFound("USER_NOT_FOUND", "rows null")
	}

	if err := p.data.db.Delete(&provider).Error; err != nil {
		return false, errors.New(500, "USER_DELETE_ERROR", "delete user error")
	}

	return true, nil
}

func (p providerRepo) GetProviderByUserId(ctx context.Context, userId uint32) (*biz.Provider, error) {
	var provider Provider
	if err := p.data.db.Where(&Provider{UserID: userId}).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
		}

		return nil, errors.New(500, "USER_NOT_FOUND", err.Error())
	}

	re := p.modelToResponse(provider)
	return re, nil
}

func (p providerRepo) GetProviderByEmail(ctx context.Context, email string) (*biz.Provider, error) {
	var provider Provider
	if err := p.data.db.Where(&Provider{Email: email}).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
		}

		return nil, errors.New(500, "USER_NOT_FOUND", err.Error())
	}

	re := p.modelToResponse(provider)
	return re, nil
}

func (p providerRepo) modelToResponse(provider Provider) *biz.Provider {
	providerInfoRsp := &biz.Provider{
		ID:           provider.ID,
		Email:        provider.Email,
		Provider:     provider.Provider,
		AccessToken:  provider.AccessToken,
		RefreshToken: provider.RefreshToken,
		TokenType:    provider.TokenType,
	}
	return providerInfoRsp
}
