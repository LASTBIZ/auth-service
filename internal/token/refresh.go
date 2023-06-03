package token

import (
	"auth-service/internal/conf"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type Refresh struct {
	t          time.Duration
	privateKey string
	publicKey  string
}

func NewRefresh(cfg *conf.Auth_Refresh) *Refresh {
	return &Refresh{
		t:          cfg.Expiry.AsDuration(),
		privateKey: cfg.Private,
		publicKey:  cfg.Public,
	}
}

func (a *Refresh) CreateToken(payload interface{}) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(a.privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(a.t).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (a *Refresh) ValidateToken(token string) (interface{}, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(a.publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims["sub"], nil
}
