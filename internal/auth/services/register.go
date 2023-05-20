package services

import (
	"context"
	"lastbiz/auth-service/internal/password"
	"lastbiz/auth-service/internal/utils"
	"lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/project"
	"lastbiz/auth-service/pkg/user"
	"net/http"
	"strings"
)

func (s Service) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if req.GetProvider() != "" {
		state, err := s.authRedis.CreateState(ctx)
		if err != nil {
			return &auth.RegisterResponse{
				Status: http.StatusInternalServerError,
				Error:  "error create state",
			}, err
		}
		providerByName, err := s.providerService.GetProviderByName(req.GetProvider())

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

		u := response.GetUser()

		createInvestor := &project.Investor{
			FullName: u.GetFirstName() + " " + u.GetLastName(),
			Money:    0,
			UserId:   u.GetId(),
		}

		projRes, err := s.projectService.AddInvestor(ctx, &project.AddInvestorRequest{
			Investor: createInvestor,
		})

		if err != nil {
			s.userService.DeleteUser(ctx, &user.DeleteUserRequest{UserId: u.GetId()})
			return &auth.RegisterResponse{
				Status: http.StatusInternalServerError,
				Error:  "error create user",
			}, err
		}

		if projRes.Status != 201 {
			s.userService.DeleteUser(ctx, &user.DeleteUserRequest{UserId: u.GetId()})
			return &auth.RegisterResponse{
				Status: response.GetStatus(),
				Error:  response.GetError(),
			}, err
		}

		hash := &password.Hash{
			Hash:   hashPass,
			UserID: u.GetId(),
		}

		err = s.passService.CreatePassword(hash)

		if err != nil {
			s.userService.DeleteUser(ctx, &user.DeleteUserRequest{UserId: u.GetId()})
			//TODO remove investor
			//s.projectService.RemoveInvestor(ctx, &project.RemoveInvestorRequest{})
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
