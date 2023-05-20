package services

import (
	"context"
	"lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/user"
	"net/http"
	"strings"
)

func (s Service) RefreshToken(ctx context.Context, request *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	token := request.GetRefreshToken()
	if strings.TrimSpace(token) == "" {
		return &auth.RefreshTokenResponse{
			Status: http.StatusNotFound,
			Error:  "token not found",
		}, nil
	}

	tok, err := s.Jwt.ValidateToken(token)
	if err != nil {
		return &auth.RefreshTokenResponse{
			Status: http.StatusNotFound,
			Error:  "token not found",
		}, nil
	}

	u, err := s.userService.GetUser(ctx, &user.UserGetRequest{
		UserId: tok.Id,
	})

	accessToken, err := s.Jwt.GenerateTokenAccess(u.GetUser())

	if err != nil {
		return &auth.RefreshTokenResponse{
			Status: http.StatusInternalServerError,
			Error:  "error create token",
		}, nil
	}

	refreshToken, err := s.Jwt.GenerateTokenRefresh(u.GetUser())

	if err != nil {
		return &auth.RefreshTokenResponse{
			Status: http.StatusInternalServerError,
			Error:  "error create token",
		}, nil
	}

	return &auth.RefreshTokenResponse{
		Status:       http.StatusOK,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
