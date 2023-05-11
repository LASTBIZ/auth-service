package provider

import (
	"context"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"lastbiz/auth-service/internal/auth"
	"testing"
	"time"
)

func TestStorage_CreateProvider(t *testing.T) {
	db, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer db.Close()
	storage := NewProviderStorage(db)

	_provider := auth.OAuthProvider{
		UserID:        1,
		OAuthProvider: "google",
		AccessToken:   "sdsadasd",
		RefreshToken:  "reasd",
		ExpiryDate:    time.Now().Add(time.Minute),
	}

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
				db.ExpectExec("INSERT INTO public.providers").
					WithArgs(
						_provider.AccessToken,
						_provider.ExpiryDate,
						_provider.OAuthProvider,
						_provider.RefreshToken,
						_provider.UserID,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			input: map[string]interface{}{
				"user_id":       _provider.UserID,
				"provider":      _provider.OAuthProvider,
				"access_token":  _provider.AccessToken,
				"refresh_token": _provider.RefreshToken,
				"expiry_date":   _provider.ExpiryDate,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := storage.CreateProvider(context.TODO(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, db.ExpectationsWereMet())
		})
	}
}

func TestStorage_UpdateProvider(t *testing.T) {
	db, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer db.Close()
	storage := NewProviderStorage(db)

	_provider := auth.OAuthProvider{
		ID:            1,
		UserID:        1,
		OAuthProvider: "google",
		AccessToken:   "sdsadasd",
		RefreshToken:  "reasd",
		ExpiryDate:    time.Now().Add(time.Minute),
	}

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
				db.ExpectExec("UPDATE public.providers").
					WithArgs(
						_provider.AccessToken,
						_provider.ExpiryDate,
						_provider.RefreshToken,
						_provider.OAuthProvider,
						_provider.UserID,
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			input: map[string]interface{}{
				"access_token":  _provider.AccessToken,
				"refresh_token": _provider.RefreshToken,
				"expiry_date":   _provider.ExpiryDate,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := storage.UpdateProvider(context.TODO(), _provider.OAuthProvider, _provider.UserID, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, db.ExpectationsWereMet())
		})
	}
}

func TestStorage_DeleteProvider(t *testing.T) {
	db, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer db.Close()
	storage := NewProviderStorage(db)

	_provider := auth.OAuthProvider{
		ID:            1,
		UserID:        1,
		OAuthProvider: "google",
		AccessToken:   "sdsadasd",
		RefreshToken:  "reasd",
		ExpiryDate:    time.Now().Add(time.Minute),
	}

	tests := []struct {
		name    string
		mock    func()
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				db.ExpectExec("DELETE FROM public.providers").
					WithArgs(_provider.OAuthProvider, _provider.UserID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := storage.DeleteProvider(context.TODO(), _provider.UserID, _provider.OAuthProvider)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, db.ExpectationsWereMet())
		})
	}
}

func TestStorage_GetProvider(t *testing.T) {
	db, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer db.Close()
	storage := NewProviderStorage(db)

	_provider := auth.OAuthProvider{
		ID:            1,
		UserID:        1,
		OAuthProvider: "google",
		AccessToken:   "sdsadasd",
		RefreshToken:  "reasd",
		ExpiryDate:    time.Now().Add(time.Minute),
	}

	tests := []struct {
		name    string
		mock    func()
		want    *auth.OAuthProvider
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rs1 := pgxmock.NewRows([]string{
					"id",
					"user_id",
					"provider",
					"access_token",
					"refresh_token",
					"expiry_date",
				}).AddRow(
					_provider.ID,
					_provider.UserID,
					_provider.OAuthProvider,
					_provider.AccessToken,
					_provider.RefreshToken,
					_provider.ExpiryDate,
				)
				db.ExpectQuery("SELECT (.+) FROM public.providers").
					WithArgs(_provider.OAuthProvider, _provider.UserID).
					WillReturnRows(rs1)
			},
			want: &_provider,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			hash, err := storage.GetProvider(context.TODO(), _provider.UserID, _provider.OAuthProvider)

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
