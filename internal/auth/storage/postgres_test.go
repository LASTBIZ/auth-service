package storage

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lastbiz/auth-service/internal/auth"
	"lastbiz/auth-service/internal/utils"
	"regexp"
	"testing"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	storage auth.Storage
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(
		postgres.Config{
			Conn: db,
		}))
	assert.Nil(s.T(), err)

	s.storage = NewAuthStorage(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_Create_Password() {
	var (
		userId   = 1
		password = "test"
	)
	hashPass := utils.HashPassword(password)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "password_hashes" ("user_id","hash") VALUES (?,?)`)).
		WithArgs(uint32(userId), hashPass).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()
	err := s.storage.CreatePassword(auth.PasswordHash{
		UserID: uint32(userId),
		Hash:   hashPass,
	})
	require.NoError(s.T(), err)
}
