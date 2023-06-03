package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
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

//go:generate mockgen -destination=../mocks/mrepo/user.go -package=mrepo . ProviderRepo
type ProviderRepo interface {
	CreateProvider(context.Context, *Provider) (*Provider, error)
	UpdateProvider(ctx context.Context, provider *Provider) (bool, error)
	DeleteProviderByUserId(ctx context.Context, userId uint32) (bool, error)
	GetProviderByUserId(ctx context.Context, userId uint32) (*Provider, error)
	GetProviderByEmail(ctx context.Context, email string) (*Provider, error)
	CreateState(ctx context.Context) (string, error)
	CheckState(ctx context.Context, state string) error
}

type ProviderUseCase struct {
	repo ProviderRepo
	log  *log.Helper
}

func NewProviderUseCase(repo ProviderRepo, logger log.Logger) *ProviderUseCase {
	return &ProviderUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (pc *ProviderUseCase) Create(ctx context.Context, p *Provider) (*Provider, error) {
	return pc.Create(ctx, p)
}

func (pc *ProviderUseCase) Update(ctx context.Context, p *Provider) (bool, error) {
	return pc.repo.UpdateProvider(ctx, p)
}

func (pc *ProviderUseCase) DeleteProviderByUserId(ctx context.Context, userId uint32) (bool, error) {
	return pc.repo.DeleteProviderByUserId(ctx, userId)
}

func (pc *ProviderUseCase) GetProviderByUserId(ctx context.Context, userId uint32) (*Provider, error) {
	return pc.repo.GetProviderByUserId(ctx, userId)
}

func (pc *ProviderUseCase) GetProviderByEmail(ctx context.Context, email string) (*Provider, error) {
	return pc.repo.GetProviderByEmail(ctx, email)
}

func (pc *ProviderUseCase) CreateState(ctx context.Context) (string, error) {
	return pc.repo.CreateState(ctx)
}

func (pc *ProviderUseCase) CheckState(ctx context.Context, state string) error {
	return pc.repo.CheckState(ctx, state)
}
