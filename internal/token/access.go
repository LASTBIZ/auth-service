package token

import (
	"auth-service/internal/conf"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/golang-jwt/jwt"
	"time"
)

type Access struct {
	t          time.Duration
	privateKey string
	publicKey  string
}

func NewAccess(cfg *conf.Auth_Access) *Access {
	return &Access{
		t:          cfg.Expiry.AsDuration(),
		privateKey: cfg.Private,
		publicKey:  cfg.Public,
	}
}

func (a *Access) CreateToken(payload interface{}) (string, error) {
	//decodedPrivateKey, err := base64.StdEncoding.DecodeString(a.privateKey)
	//if err != nil {
	//	return "", fmt.Errorf("could not decode key: %w", err)
	//}
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(a.privateKey))

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

func (a *Access) ValidateToken(token string) (interface{}, error) {
	//decodedPublicKey, err := base64.StdEncoding.DecodeString(a.publicKey)
	//if err != nil {
	//	return nil, fmt.Errorf("could not decode: %w", err)
	//}

	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(a.publicKey))

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
	if !ok {
		return nil, fmt.Errorf("validate: invalid token")
	}

	err = claims.Valid()
	if err != nil {
		if err.Error() == "Token is expired" {
			return nil, errors.InternalServer("EXPIRY_TOKEN", "token expiry")
		}
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims["sub"], nil
}
