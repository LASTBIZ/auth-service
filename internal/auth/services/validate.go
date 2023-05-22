package services

import (
	"context"
	"lastbiz/auth-service/internal/utils"
	"lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/user"
	"net/http"
	"strings"
)

func (s Service) Validate(ctx context.Context, request *auth.ValidateRequest) (*auth.ValidateResponse, error) {
	token := request.GetToken()
	if strings.TrimSpace(token) == "" {
		return &auth.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "token not found",
		}, nil
	}

	tok, err := utils.ValidateToken(token, s.PrivateKeyAccess)
	if err != nil {
		return &auth.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "token not found",
		}, nil
	}
	//check user is exists
	u, err := s.userService.GetUser(ctx, &user.UserGetRequest{
		UserId: tok.(uint32),
	})
	if u.Status != 200 {
		return &auth.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "user not found",
		}, nil
	}

	return &auth.ValidateResponse{
		Status: http.StatusOK,
		UserId: int64(tok.(uint32)),
	}, nil
}
