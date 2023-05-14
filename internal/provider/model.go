package provider

import (
	"github.com/devfeel/mapper"
	"lastbiz/auth-service/pkg/errors"
	"time"
)

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
	return providerMap, nil
}
