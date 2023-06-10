package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type Session struct {
	ID           uint64
	UserID       uint64
	RefreshToken string
	//ua           string
	//ip           string
	ExpiresIn time.Time
	CreatedAt time.Time
}

//go:generate mockgen -destination=../mocks/mrepo/user.go -package=mrepo . SessionRepo
type SessionRepo interface {
	CreateSession(context.Context, *Session) (*Session, error)
	UpdateSession(ctx context.Context, session *Session) (bool, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*Session, error)
	GetSessionByUserID(ctx context.Context, userId uint64) (*Session, error)
}

type SessionUseCase struct {
	repo SessionRepo
	log  *log.Helper
}

func NewSessionUseCase(repo SessionRepo, logger log.Logger) *SessionUseCase {
	return &SessionUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (hc *SessionUseCase) CreateSession(ctx context.Context, session *Session) (*Session, error) {
	return hc.repo.CreateSession(ctx, session)
}

func (hc *SessionUseCase) UpdateSession(ctx context.Context, session *Session) (bool, error) {
	return hc.repo.UpdateSession(ctx, session)
}

func (hc *SessionUseCase) GetSession(ctx context.Context, refreshToken string) (*Session, error) {
	return hc.repo.GetSessionByRefreshToken(ctx, refreshToken)
}

func (hc *SessionUseCase) GetSessionByUserID(ctx context.Context, userId uint64) (*Session, error) {
	return hc.repo.GetSessionByUserID(ctx, userId)
}
