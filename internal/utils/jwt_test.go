package utils

import (
	"github.com/stretchr/testify/assert"
	"lastbiz/auth-service/pkg/user"
	"testing"
)

func TestHashPassword(t *testing.T) {

	jwt := &JwtWrapper{
		SecretKey:       "",
		Issuer:          "",
		ExpirationHours: 100,
	}

	tests := []struct {
		name    string
		mock    string
		input   *user.User
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			input: &user.User{
				Id:    1,
				Email: "test@gmail.com",
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := jwt.GenerateToken(tt.input)
			claims, err := jwt.ValidateToken(token)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Id, claims.Id)
				assert.Equal(t, tt.input.Email, claims.Email)
			}
		})
	}
}
