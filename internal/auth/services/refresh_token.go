package services

import (
	"context"
	"lastbiz/auth-service/internal/utils"
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

	tok, err := utils.ValidateToken(token, s.PrivateKeyRefresh)
	if err != nil {
		return &auth.RefreshTokenResponse{
			Status: http.StatusNotFound,
			Error:  "token not found",
		}, nil
	}

	u, err := s.userService.GetUser(ctx, &user.UserGetRequest{
		UserId: tok.(uint32),
	})

	if err != nil {
		return &auth.RefreshTokenResponse{
			Status: http.StatusInternalServerError,
			Error:  "the user belonging to this token no logger exists",
		}, nil
	}

	if u.GetStatus() != 200 {
		return &auth.RefreshTokenResponse{
			Status: http.StatusInternalServerError,
			Error:  "the user belonging to this token no logger exists",
		}, nil
	}

	accessToken, err := utils.CreateToken(s.AccessTokenDuration, u.GetUser().GetId(), s.PrivateKeyAccess)

	if err != nil {
		return &auth.RefreshTokenResponse{
			Status: http.StatusInternalServerError,
			Error:  "error create token",
		}, nil
	}

	refreshToken, err := utils.CreateToken(s.RefreshTokenDuration, u.GetUser().GetId(), s.PrivateKeyRefresh)

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
