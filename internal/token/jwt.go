package token

import "auth-service/internal/conf"

type JwtClaims struct {
	Access  *Access
	Refresh *Refresh
}

func NewJwtClaims(cfg *conf.Auth) *JwtClaims {
	access := NewAccess(cfg.Access)
	refresh := NewRefresh(cfg.Refresh)
	return &JwtClaims{
		access,
		refresh,
	}
}
