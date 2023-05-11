package password

import (
	"context"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"lastbiz/auth-service/internal/auth"
	"lastbiz/auth-service/internal/utils"
	"testing"
)

func Test_Create_Password(t *testing.T) {
	db, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer db.Close()
	storage := NewPasswordStorage(db)
	var (
		userId   = 1
		password = "test"
	)
	hashPass := utils.HashPassword(password)
	tests := []struct {
		name    string
		mock    func()
		input   map[string]interface{}
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				db.ExpectExec("INSERT INTO public.passwords").
					WithArgs(hashPass, userId).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			input: map[string]interface{}{
				"user_id": userId,
				"hash":    hashPass,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := storage.CreatePassword(context.TODO(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, db.ExpectationsWereMet())
		})
	}
}

func TestStorage_UpdatePassword(t *testing.T) {
	db, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer db.Close()
	storage := NewPasswordStorage(db)
	var (
		id       = 1
		userId   = 1
		password = "test"
	)
	hashPass := utils.HashPassword(password)
	tests := []struct {
		name    string
		mock    func()
		input   map[string]interface{}
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				db.ExpectExec("UPDATE public.passwords").
					WithArgs(hashPass, userId, uint32(id)).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			input: map[string]interface{}{
				"user_id": userId,
				"hash":    hashPass,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := storage.UpdatePassword(context.TODO(), uint32(1), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, db.ExpectationsWereMet())
		})
	}
}

func TestStorage_DeletePassword(t *testing.T) {
	db, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer db.Close()
	storage := NewPasswordStorage(db)
	var (
		id = uint32(1)
	)
	tests := []struct {
		name    string
		mock    func()
		input   uint32
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				db.ExpectExec("DELETE FROM public.passwords").
					WithArgs(id).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			input: id,
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := storage.DeletePassword(context.TODO(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, db.ExpectationsWereMet())
		})
	}
}

func TestStorage_GetHash(t *testing.T) {
	db, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer db.Close()
	storage := NewPasswordStorage(db)
	var (
		id     = uint32(1)
		hash   = "asdasdasdasdasdasd"
		userID = uint32(1)
	)
	tests := []struct {
		name    string
		mock    func()
		input   uint32
		want    *auth.PasswordHash
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rs1 := pgxmock.NewRows([]string{"id", "user_id", "hash"}).AddRow(id, userID, hash)
				db.ExpectQuery("SELECT (.+) FROM public.passwords WHERE user_id = ?").
					WithArgs(userID).
					WillReturnRows(rs1)
			},
			input: userID,
			want:  &auth.PasswordHash{ID: id, Hash: hash, UserID: userID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			hash, err := storage.GetHash(context.TODO(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, hash)
			}
			assert.NoError(t, db.ExpectationsWereMet())
		})
	}
}
