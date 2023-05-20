package services

import (
	"context"
	"lastbiz/auth-service/internal/provider"
	"lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/user"
	"net/http"
	"strings"
)

func (s Service) Callback(ctx context.Context, request *auth.CallbackRequest) (*auth.CallbackResponse, error) {
	state := request.GetState()
	if strings.TrimSpace(state) == "" {
		return &auth.CallbackResponse{
			Status: http.StatusConflict,
			Error:  "state not found",
		}, nil
	}

	ok, err := s.authRedis.CheckState(ctx, state)
	if err != nil {
		return &auth.CallbackResponse{
			Status: http.StatusConflict,
			Error:  "state not found",
		}, nil
	}
	if !ok {
		return &auth.CallbackResponse{
			Status: http.StatusConflict,
			Error:  "state not found",
		}, nil
	}

	providerName := request.GetProvider()

	pr, err := s.providerService.GetProviderByName(providerName)
	if err != nil {
		return &auth.CallbackResponse{
			Status: http.StatusConflict,
			Error:  "provider not found",
		}, nil
	}

	tokenSource, err := pr.Callback(request.GetOauthCode())
	if err != nil {
		return &auth.CallbackResponse{
			Status: http.StatusInternalServerError,
			Error:  "error validate code",
		}, nil
	}

	providerUser, err := pr.GetUser(tokenSource)
	if err != nil {
		return &auth.CallbackResponse{
			Status: http.StatusInternalServerError,
			Error:  "error getUser",
		}, nil
	}

	result, err := s.userService.GetUserByEmail(ctx, &user.UserByEmailRequest{
		Email: providerUser.Email,
	})
	u := &user.User{}

	if result.Status == 404 && result.GetError() == "User not found" {
		//Create user and provider
		createUser := &user.User{
			LastName:  providerUser.GivenName,
			FirstName: providerUser.FamilyName,
			Email:     providerUser.Email,
			IsVerify:  providerUser.VerifyEmail,
		}
		resultUser, err := s.userService.CreateUser(ctx, createUser)
		if err != nil {
			return &auth.CallbackResponse{
				Status: http.StatusInternalServerError,
				Error:  "error create user",
			}, nil
		}
		if resultUser.Status != 201 {
			//todo catch error
			return &auth.CallbackResponse{
				Status: http.StatusInternalServerError,
				Error:  "error create user",
			}, nil
		}

		createProvider := &provider.OAuthProvider{
			UserID:       resultUser.GetUser().GetId(),
			Provider:     strings.ToLower(providerName),
			AccessToken:  tokenSource.AccessToken,
			RefreshToken: tokenSource.RefreshToken,
			ExpiryDate:   tokenSource.Expiry,
			TokenType:    tokenSource.TokenType,
		}
		err = s.providerService.CreateProvider(createProvider)
		if err != nil {
			return &auth.CallbackResponse{
				Status: http.StatusInternalServerError,
				Error:  "error create provider",
			}, nil
		}

		u = resultUser.User
	}

	if result.GetStatus() == 200 {
		resultUser, err := s.userService.GetUserByEmail(ctx, &user.UserByEmailRequest{
			Email: providerUser.Email,
		})
		if err != nil {
			return &auth.CallbackResponse{
				Status: http.StatusInternalServerError,
				Error:  "error get user",
			}, nil
		}
		if resultUser.GetStatus() != 200 {
			return &auth.CallbackResponse{
				Status: http.StatusInternalServerError,
				Error:  "error get user",
			}, nil
		}
		isRegister, err := s.CheckRegister(ctx, &auth.CheckRegisterRequest{
			UserId: int64(resultUser.GetUser().GetId()),
			Pass:   false,
		})
		if isRegister.GetStatus() != 200 {
			return &auth.CallbackResponse{
				Status: http.StatusInternalServerError,
				Error:  "error get user",
			}, nil
		}
		if isRegister.GetRegistered() {
			err = s.providerService.UpdateProvider(
				strings.ToLower(providerName), resultUser.GetUser().GetId(), tokenSource.AccessToken, tokenSource.RefreshToken, tokenSource.Expiry)
			return &auth.CallbackResponse{
				Status: http.StatusInternalServerError,
				Error:  "error update provider",
			}, nil
		} else {
			createProvider := &provider.OAuthProvider{
				UserID:       resultUser.GetUser().GetId(),
				Provider:     strings.ToLower(providerName),
				AccessToken:  tokenSource.AccessToken,
				RefreshToken: tokenSource.RefreshToken,
				ExpiryDate:   tokenSource.Expiry,
				TokenType:    tokenSource.TokenType,
			}
			err = s.providerService.CreateProvider(createProvider)
			if err != nil {
				return &auth.CallbackResponse{
					Status: http.StatusInternalServerError,
					Error:  "error create provider",
				}, nil
			}
		}

		u = resultUser.GetUser()
	} else {
		return &auth.CallbackResponse{
			Status: http.StatusNotFound,
			Error:  "not found",
		}, nil
	}

	//User exists login
	//generate token access_token refresh_token
	accessToken, err := s.Jwt.GenerateTokenAccess(u)
	refreshToken, err := s.Jwt.GenerateTokenRefresh(u)
	return &auth.CallbackResponse{
		Status: http.StatusOK,
		Token: &auth.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
