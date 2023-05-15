package auth

import (
	"context"
	"lastbiz/auth-service/internal/password"
	"lastbiz/auth-service/internal/provider"
	"lastbiz/auth-service/internal/utils"
	"lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/user"
	"net/http"
	"strings"
)

type Service struct {
	passService     password.Service
	providerService provider.Service
	Jwt             utils.JwtWrapper
	userService     user.UserServiceClient
}

func NewAuthService(
	passService password.Service,
	providerService provider.Service,
	userService user.UserServiceClient,
) *Service {
	return &Service{
		passService:     passService,
		providerService: providerService,
		userService:     userService,
	}
}

func (s Service) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if req.GetCode() != nil {
		//TODO add redis for state save and check
		state := utils.GenerateState()
		providerByName, err := s.providerService.GetProviderByName(req.GetCode().GetProvider())

		if err != nil {
			return &auth.RegisterResponse{
				Status: http.StatusNotFound,
				Error:  "providerByName not found",
			}, err
		}

		url := providerByName.GenerateOAuthToken(state)

		return &auth.RegisterResponse{
			Status:   http.StatusPermanentRedirect,
			Redirect: url,
		}, err
	} else if strings.TrimSpace(req.GetPassword()) != "" {
		hashPass := utils.HashPassword(req.GetPassword())

		createUser := &user.User{
			Email:     req.GetEmail(),
			FirstName: req.GetFirstName(),
			LastName:  req.GetLastName(),
		}

		response, err := s.userService.CreateUser(ctx, createUser)

		if err != nil {
			return &auth.RegisterResponse{
				Status: http.StatusInternalServerError,
				Error:  "error create user",
			}, err
		}

		if response.Status != 201 {
			return &auth.RegisterResponse{
				Status: response.GetStatus(),
				Error:  response.GetError(),
			}, err
		}

		hash := &password.Hash{
			Hash:   hashPass,
			UserID: response.GetUser().GetId(),
		}

		err = s.passService.CreatePassword(hash)

		if err != nil {
			return &auth.RegisterResponse{
				Status: http.StatusInternalServerError,
				Error:  "error register",
			}, err
		}

		return &auth.RegisterResponse{
			Status: http.StatusCreated,
		}, err

	} else {
		return &auth.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "error method register",
		}, nil
	}
}

func (s Service) Login(ctx context.Context) {

}

func (s Service) LoginByProvider(ctx context.Context) {

}
