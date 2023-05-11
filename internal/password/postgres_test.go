package password

import (
	"context"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"lastbiz/auth-service/internal/utils"
	"testing"
)

type Suite struct {
	suite.Suite
	mock pgxmock.PgxPoolIface

	storage *Storage
}

func Test_Create_Password(t *testing.T) {
	db, err := pgxmock.NewPool()
	assert.NoError(t, err)

	storage := NewPasswordStorage(db)
	var (
		userId   = 1
		password = "test"
	)

	hashPass := utils.HashPassword(password)
	db.ExpectExec("INSERT INTO passwords").
		WithArgs(hashPass, userId).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = storage.CreatePassword(context.TODO(), map[string]interface{}{
		"user_id": userId,
		"hash":    hashPass,
	})

	assert.NoError(t, err)
}
