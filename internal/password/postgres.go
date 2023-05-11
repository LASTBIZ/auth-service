package password

import (
	"context"
	"errors"
	"fmt"
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

func NewPasswordStorage(client storage2.PostgreSQLClient) *Storage {
	return &Storage{
		client:       client,
		queryBuilder: sq.StatementBuilder,
	}
}

const (
	scheme      = "public"
	table       = "passwords"
	tableScheme = scheme + "." + table
)

func (s Storage) CreatePassword(ctx context.Context, m map[string]interface{}) error {
	sql, args, buildErr := s.queryBuilder.
		Insert(table).
		SetMap(m).
		ToSql()
	logger := logging.WithFields(ctx, map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
	})

	fmt.Println(sql)

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
		execErr = db.ErrDoQuery(errors.New("password was not created. 0 rows were affected"))
		logger.Error(execErr)
		return execErr
	}

	return nil
}

func (s Storage) UpdatePassword(ctx context.Context, id uint32, m map[string]interface{}) error {
	sql, args, buildErr := s.queryBuilder.
		Update(tableScheme).
		SetMap(m).
		Where(sq.Eq{"id": id}).
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
		execErr = db.ErrDoQuery(errors.New("password was not created. 0 rows were affected"))
		logger.Error(execErr)
		return execErr
	}

	return nil
}

func (s Storage) DeletePassword(ctx context.Context, id uint32) error {
	sql, args, buildErr := s.queryBuilder.
		Delete(tableScheme).
		Where(sq.Eq{"id": id}).
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
		execErr = db.ErrDoQuery(errors.New("password was not created. 0 rows were affected"))
		logger.Error(execErr)
		return execErr
	}

	return nil
}

func (s Storage) GetHash(ctx context.Context, userID uint32) (*auth.PasswordHash, error) {
	sql, args, buildErr := s.queryBuilder.
		Select("id").
		Columns(
			"user_id",
			"hash").
		From(tableScheme).
		Where(sq.Eq{"user_id": userID}).
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

	var hash auth.PasswordHash

	err := s.client.QueryRow(ctx, sql, args...).Scan(
		&hash.ID,
		&hash.UserID,
		&hash.Hash,
	)
	if err != nil {
		err = db.ErrDoQuery(err)
		logger.Error(err)
		return nil, err
	}

	return &hash, err
}
