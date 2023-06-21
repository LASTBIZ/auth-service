package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type Hash struct {
	ID     uint32
	UserID uint32
	Hash   string
}

//go:generate mockgen -destination=../mocks/mrepo/hash.go -package=mrepo . HashRepo
type HashRepo interface {
	CreateHash(context.Context, *Hash) (*Hash, error)
	UpdateHash(ctx context.Context, hash *Hash) (bool, error)
	DeleteHashByUserId(ctx context.Context, userId uint32) (bool, error)
	GetHashByUserId(ctx context.Context, userId uint32) (*Hash, error)
}

type HashUseCase struct {
	repo HashRepo
	log  *log.Helper
}

func NewHashUseCase(repo HashRepo, logger log.Logger) *HashUseCase {
	return &HashUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (hc *HashUseCase) Create(ctx context.Context, h *Hash) (*Hash, error) {
	return hc.repo.CreateHash(ctx, h)
}

func (hc *HashUseCase) Update(ctx context.Context, h *Hash) (bool, error) {
	return hc.repo.UpdateHash(ctx, h)
}

func (hc *HashUseCase) DeleteHashByUserId(ctx context.Context, userId uint32) (bool, error) {
	return hc.repo.DeleteHashByUserId(ctx, userId)
}

func (hc *HashUseCase) GetHashByUserId(ctx context.Context, userId uint32) (*Hash, error) {
	return hc.repo.GetHashByUserId(ctx, userId)
}
