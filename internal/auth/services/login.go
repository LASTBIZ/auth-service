package services

import (
	"context"
	"lastbiz/auth-service/internal/utils"
	"lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/user"
	"net/http"
	"strings"
)

func (s Service) LoginByPassword(ctx context.Context, request *auth.LoginByPasswordRequest) (*auth.LoginResponse, error) {
	if strings.TrimSpace(request.GetPassword()) == "" {
		return &auth.LoginResponse{
			Status: http.StatusConflict,
			Error:  "password is empty",
		}, nil
	}

	if strings.TrimSpace(request.GetEmail()) == "" {
		return &auth.LoginResponse{
			Status: http.StatusConflict,
			Error:  "email is empty",
		}, nil
	}

	u, err := s.userService.GetUserByEmail(ctx, &user.UserByEmailRequest{Email: request.GetEmail()})
	if err != nil {
		return &auth.LoginResponse{
			Status: http.StatusConflict,
			Error:  "wrong password or user not found",
		}, nil
	}

	if u.Status != http.StatusOK {
		return &auth.LoginResponse{
			Status: http.StatusConflict,
			Error:  "wrong password or user not found",
		}, nil
	}

	hash, err := s.passService.GetHash(u.GetUser().GetId())
	if err != nil {
		return &auth.LoginResponse{
			Status: http.StatusConflict,
			Error:  "wrong password or user not found",
		}, nil
	}
	if !utils.CheckPasswordHash(request.GetPassword(), hash.Hash) {
		return &auth.LoginResponse{
			Status: http.StatusConflict,
			Error:  "wrong password or user not found",
		}, nil
	}

	token, err := utils.CreateToken(s.AccessTokenDuration, u.GetUser().GetId(), s.PrivateKeyAccess)
	refresh, err := utils.CreateToken(s.RefreshTokenDuration, u.GetUser().GetId(), s.PrivateKeyRefresh)

	if err != nil {
		return &auth.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  "server error",
		}, nil
	}

	return &auth.LoginResponse{
		Status: http.StatusOK,
		Token: &auth.Token{
			AccessToken:  token,
			RefreshToken: refresh,
		},
	}, nil
}
