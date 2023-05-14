package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPasswordHash_ToMap(t *testing.T) {
	tests := []struct {
		name    string
		input   *PasswordHash
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Ok",
			input: &PasswordHash{
				ID:     1,
				UserID: 1,
				Hash:   "test",
			},
			want: map[string]interface{}{
				"id":      uint32(1),
				"user_id": uint32(1),
				"hash":    "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tomap, err := tt.input.ToMap()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tomap, tt.want)
			}
		})
	}
}

func TestOAuthProvider_ToMap(t *testing.T) {
	time := time.Now()
	tests := []struct {
		name    string
		input   *OAuthProvider
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Ok",
			input: &OAuthProvider{
				ID:            1,
				UserID:        1,
				OAuthProvider: "google",
				AccessToken:   "test",
				RefreshToken:  "test",
				ExpiryDate:    time,
			},
			want: map[string]interface{}{
				"id":            uint32(1),
				"user_id":       uint32(1),
				"provider":      "google",
				"access_token":  "test",
				"refresh_token": "test",
				"expiry_date":   time,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tomap, err := tt.input.ToMap()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tomap, tt.want)
			}
		})
	}
}
