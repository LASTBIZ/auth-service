package services

import (
	"context"
	"gorm.io/gorm"
	"lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/errors"
	"lastbiz/auth-service/pkg/user"
	"net/http"
)

func (s Service) CheckRegister(ctx context.Context, req *auth.CheckRegisterRequest) (*auth.CheckRegisterResponse, error) {
	userId := uint32(0)
	if req.GetEmail() != "" {
		resultUser, err := s.userService.GetUserByEmail(ctx, &user.UserByEmailRequest{Email: req.GetEmail()})
		if err != nil {
			return &auth.CheckRegisterResponse{
				Status: http.StatusInternalServerError,
				Error:  "error get user",
			}, nil
		}

		if resultUser.GetStatus() != 200 {
			return &auth.CheckRegisterResponse{
				Status: http.StatusInternalServerError,
				Error:  "error get user",
			}, nil
		}
		userId = resultUser.GetUser().GetId()
	} else if req.GetUserId() == 0 {
		return &auth.CheckRegisterResponse{
			Status: http.StatusInternalServerError,
			Error:  "error get user",
		}, nil
	}
	isRegister := true

	if req.GetPass() {
		_, err := s.passService.GetHash(userId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				isRegister = false
			}
			return &auth.CheckRegisterResponse{
				Status: http.StatusInternalServerError,
				Error:  "error get hash",
			}, nil
		} else {
			isRegister = true
		}
	} else {
		_, err := s.providerService.CheckProvider(userId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				isRegister = false
			}
			return &auth.CheckRegisterResponse{
				Status: http.StatusInternalServerError,
				Error:  "error get provider",
			}, nil
		} else {
			isRegister = true
		}
	}

	return &auth.CheckRegisterResponse{
		Status:     http.StatusOK,
		Registered: isRegister,
	}, nil
}
