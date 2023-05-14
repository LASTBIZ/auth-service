package auth

import (
	"fmt"
	"github.com/devfeel/mapper"
	"lastbiz/auth-service/pkg/errors"
	"time"
)

type PasswordHash struct {
	ID     uint32 `mapper:"id"`
	UserID uint32 `mapper:"user_id"`
	Hash   string `mapper:"hash"`
}

func (p *PasswordHash) ToMap() (map[string]interface{}, error) {
	hashMap := make(map[string]interface{})
	err := mapper.AutoMapper(p, &hashMap)
	if err != nil {
		return hashMap, errors.Wrap(err, "mapper.Decode(password_hash)")
	}

	return hashMap, nil
}

type OAuthProvider struct {
	ID            uint32    `mapper:"id"`
	UserID        uint32    `mapper:"user_id"`
	OAuthProvider string    `mapper:"provider"`
	AccessToken   string    `mapper:"access_token"`
	RefreshToken  string    `mapper:"refresh_token"`
	ExpiryDate    time.Time `mapper:"expiry_date"`
}

func (p *OAuthProvider) ToMap() (map[string]interface{}, error) {
	providerMap := make(map[string]interface{})
	err := mapper.AutoMapper(p, &providerMap)
	if err != nil {
		return providerMap, errors.Wrap(err, "mapper.Decode(provider)")
	}
	fmt.Println(providerMap)
	return providerMap, nil
}
