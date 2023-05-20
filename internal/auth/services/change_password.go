package services

import (
	"context"
	"lastbiz/auth-service/internal/utils"
	"lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/user"
	"net/http"
	"strings"
)

func (s Service) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	if strings.TrimSpace(req.GetPassword()) == "" {
		return &auth.ChangePasswordResponse{
			Status: http.StatusConflict,
			Error:  "password is empty",
		}, nil
	}

	if strings.TrimSpace(req.GetNewPassword()) == "" {
		return &auth.ChangePasswordResponse{
			Status: http.StatusConflict,
			Error:  "password is empty",
		}, nil
	}

	if strings.TrimSpace(req.GetEmail()) == "" {
		return &auth.ChangePasswordResponse{
			Status: http.StatusConflict,
			Error:  "password is empty",
		}, nil
	}

	//get user
	u, err := s.userService.GetUserByEmail(ctx, &user.UserByEmailRequest{Email: req.GetEmail()})
	if err != nil {
		return &auth.ChangePasswordResponse{
			Status: http.StatusInternalServerError,
			Error:  "error get user",
		}, nil
	}

	if u.GetStatus() != 200 {
		return &auth.ChangePasswordResponse{
			Status: http.StatusConflict,
			Error:  "user not found",
		}, nil
	}

	hash, err := s.passService.GetHash(u.GetUser().GetId())
	if err != nil {
		return &auth.ChangePasswordResponse{
			Status: http.StatusInternalServerError,
			Error:  "error get user",
		}, nil
	}

	if !utils.CheckPasswordHash(req.GetPassword(), req.GetNewPassword()) {
		return &auth.ChangePasswordResponse{
			Status: http.StatusInternalServerError,
			Error:  "wrong password",
		}, nil
	}

	hash.Hash = utils.HashPassword(req.GetNewPassword())

	err = s.passService.UpdatePassword(u.GetUser().GetId(), hash)
	if err != nil {
		return &auth.ChangePasswordResponse{
			Status: http.StatusInternalServerError,
			Error:  "error update password",
		}, nil
	}
	return &auth.ChangePasswordResponse{
		Status: http.StatusOK,
	}, nil
}
