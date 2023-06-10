package data

import (
	"auth-service/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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

type sessionRepo struct {
	data *Data
	log  *log.Helper
}

func NewSessionRepo(data *Data, logger log.Logger) biz.SessionRepo {
	return &sessionRepo{data: data, log: log.NewHelper(logger)}
}

func (s sessionRepo) CreateSession(ctx context.Context, session *biz.Session) (*biz.Session, error) {
	var sess Session
	result := s.data.db.Where(&Session{RefreshToken: session.RefreshToken}).Or(&Session{UserID: session.UserID}).First(&sess)
	if result.RowsAffected == 1 {
		return nil, errors.New(500, "SESSION_EXISTS", "session is exists")
	}

	sess.UserID = session.UserID
	sess.RefreshToken = session.RefreshToken
	sess.CreatedAt = session.CreatedAt
	sess.ExpiresIn = session.ExpiresIn

	res := s.data.db.Create(&sess)
	if res.Error != nil {
		return nil, errors.New(500, "CREATE_SESSION_ERROR", "error create session")
	}

	sessionInfoRes := s.modelToResponse(sess)
	return sessionInfoRes, nil
}

func (s sessionRepo) UpdateSession(ctx context.Context, session *biz.Session) (bool, error) {
	var sessionInfo Session
	result := s.data.db.Where(&Session{UserID: session.UserID}).First(&sessionInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, errors.NotFound("SESSION_NOT_FOUND", "session not found")
	}

	if result.RowsAffected == 0 {
		return false, errors.NotFound("SESSION_NOT_FOUND", "rows null")
	}

	sessionInfo.RefreshToken = session.RefreshToken
	sessionInfo.CreatedAt = session.CreatedAt
	sessionInfo.ExpiresIn = session.ExpiresIn

	if err := s.data.db.Save(&sessionInfo).Error; err != nil {
		return false, errors.New(500, "SESSION_UPDATE_ERROR", "update session error")
	}

	return true, nil
}

func (s sessionRepo) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*biz.Session, error) {
	var sessionInfo Session
	if err := s.data.db.Where(&Session{RefreshToken: refreshToken}).First(&sessionInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("SESSION_NOT_FOUND", "session not found")
		}

		return nil, errors.New(500, "SESSION_NOT_FOUND", err.Error())
	}

	re := s.modelToResponse(sessionInfo)
	return re, nil
}

func (s sessionRepo) GetSessionByUserID(ctx context.Context, userId uint64) (*biz.Session, error) {
	var sessionInfo Session
	if err := s.data.db.Where(&Session{UserID: userId}).First(&sessionInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("SESSION_NOT_FOUND", "session not found")
		}

		return nil, errors.New(500, "SESSION_NOT_FOUND", err.Error())
	}

	re := s.modelToResponse(sessionInfo)
	return re, nil
}

func (s sessionRepo) modelToResponse(ses Session) *biz.Session {
	sessionInfoRsp := &biz.Session{
		ID:           ses.ID,
		UserID:       ses.UserID,
		RefreshToken: ses.RefreshToken,
		CreatedAt:    ses.CreatedAt,
		ExpiresIn:    ses.ExpiresIn,
	}
	return sessionInfoRsp
}
