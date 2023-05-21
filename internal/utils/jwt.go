package utils

import (
	"errors"
	"lastbiz/auth-service/pkg/user"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtWrapper struct {
	SecretKey              string
	Issuer                 string
	ExpirationHoursAccess  int64
	ExpirationHoursRefresh int64
}

type jwtClaims struct {
	jwt.StandardClaims
	Id    uint32
	Email string
}

// TODO maybe save in redis
func (w *JwtWrapper) GenerateTokenRefresh(user *user.User) (signedToken string, err error) {
	return w.generateToken(user, w.ExpirationHoursRefresh, w.Issuer)
}

func (w *JwtWrapper) GenerateTokenAccess(user *user.User) (signedToken string, err error) {
	return w.generateToken(user, w.ExpirationHoursAccess, "refresh")
}

func (w *JwtWrapper) generateToken(user *user.User, expiryTime int64, issuer string) (signedToken string, err error) {
	claims := &jwtClaims{
		Id:    user.Id,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(expiryTime)).Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(w.SecretKey))

	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (w *JwtWrapper) ValidateToken(signedToken, tokenType string) (claims *jwtClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(w.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*jwtClaims)

	if !ok {
		return nil, errors.New("Couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("JWT is expired")
	}

	if tokenType == "refresh" && claims.Issuer != "refresh" {
		return nil, errors.New("token not refresh")
	}

	return claims, nil

}
