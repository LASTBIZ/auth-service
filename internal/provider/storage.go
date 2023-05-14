package provider

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"lastbiz/auth-service/internal/auth"
	storage2 "lastbiz/auth-service/internal/storage"
	"lastbiz/auth-service/pkg/logging"
	db "lastbiz/auth-service/pkg/postgres/model"
)

type Storage struct {
	queryBuilder sq.StatementBuilderType
	client       storage2.PostgreSQLClient
}

func NewProviderStorage(client storage2.PostgreSQLClient) *Storage {
	return &Storage{
		client:       client,
		queryBuilder: sq.StatementBuilder,
	}
}

const (
	scheme      = "public"
	table       = "providers"
	tableScheme = scheme + "." + table
)

func (s Storage) CreateProvider(ctx context.Context, m map[string]interface{}) error {
	sql, args, buildErr := s.queryBuilder.
		Insert(tableScheme).
		SetMap(m).
		ToSql()

	logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": tableScheme,
		"args":  args,
	})

	if buildErr != nil {
		buildErr = db.ErrCreateQuery(buildErr)
		logger.Error(buildErr)
		return buildErr
	}
	if exec, execErr := s.client.Exec(ctx, sql, args...); execErr != nil {
		execErr = db.ErrDoQuery(execErr)
		logger.Error(execErr)
		return execErr
	} else if exec.RowsAffected() == 0 || !exec.Insert() {
		execErr = db.ErrDoQuery(errors.New("provider was not created. 0 rows were affected"))
		logger.Error(execErr)
		return execErr
	}

	return nil
}

func (s Storage) UpdateProvider(ctx context.Context, provider string, userId uint32, m map[string]interface{}) error {
	sql, args, buildErr := s.queryBuilder.
		Update(tableScheme).
		SetMap(m).
		Where(sq.Eq{"provider": provider, "user_id": userId}).
		ToSql()

	logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": tableScheme,
		"args":  args,
	})

	if buildErr != nil {
		buildErr = db.ErrCreateQuery(buildErr)
		logger.Error(buildErr)
		return buildErr
	}

	if exec, execErr := s.client.Exec(ctx, sql, args...); execErr != nil {
		execErr = db.ErrDoQuery(execErr)
		logger.Error(execErr)
		return execErr
	} else if exec.RowsAffected() == 0 || !exec.Update() {
		execErr = db.ErrDoQuery(errors.New("password was not created. 0 rows were affected"))
		logger.Error(execErr)
		return execErr
	}

	return nil
}

func (s Storage) DeleteProvider(ctx context.Context, userID uint32, provider string) error {
	sql, args, buildErr := s.queryBuilder.
		Delete(tableScheme).
		Where(sq.Eq{"user_id": userID, "provider": provider}).
		ToSql()

	logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": tableScheme,
		"args":  args,
	})

	if buildErr != nil {
		buildErr = db.ErrCreateQuery(buildErr)
		logger.Error(buildErr)
		return buildErr
	}

	if exec, execErr := s.client.Exec(ctx, sql, args...); execErr != nil {
		execErr = db.ErrDoQuery(execErr)
		logger.Error(execErr)
		return execErr
	} else if exec.RowsAffected() == 0 || !exec.Delete() {
		execErr = db.ErrDoQuery(errors.New("password was not created. 0 rows were affected"))
		logger.Error(execErr)
		return execErr
	}

	return nil
}

func (s Storage) GetProvider(ctx context.Context, userID uint32, provider string) (*auth.OAuthProvider, error) {
	sql, args, buildErr := s.queryBuilder.
		Select("id").
		Columns(
			"user_id",
			"provider",
			"access_token",
			"refresh_token",
			"expiry_date").
		From(tableScheme).
		Where(sq.Eq{"user_id": userID, "provider": provider}).
		ToSql()
	logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": tableScheme,
		"args":  args,
	})

	if buildErr != nil {
		buildErr = db.ErrCreateQuery(buildErr)
		logger.Error(buildErr)
		return nil, buildErr
	}

	var _provider auth.OAuthProvider

	err := s.client.QueryRow(ctx, sql, args...).Scan(
		&_provider.ID,
		&_provider.UserID,
		&_provider.OAuthProvider,
		&_provider.AccessToken,
		&_provider.RefreshToken,
		&_provider.ExpiryDate,
	)

	if err != nil {
		err = db.ErrDoQuery(err)
		logger.Error(err)
		return nil, err
	}

	return &_provider, err
}
