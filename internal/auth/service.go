package auth

import (
	"context"
	"gorm.io/gorm"
	"lastbiz/auth-service/internal/password"
	"lastbiz/auth-service/internal/provider"
	"lastbiz/auth-service/internal/utils"
	"lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/errors"
	"lastbiz/auth-service/pkg/user"
	"net/http"
	"strings"
)

type Service struct {
	passService     password.Service
	providerService provider.Service
	Jwt             utils.JwtWrapper
	userService     user.UserServiceClient
	authRedis       authRedis
	auth.UnimplementedAuthServiceServer
}

func NewAuthService(
	passService password.Service,
	providerService provider.Service,
	userService user.UserServiceClient,
	authRedis authRedis,
	jwt utils.JwtWrapper,
) auth.AuthServiceServer {
	return Service{
		passService:     passService,
		providerService: providerService,
		userService:     userService,
		authRedis:       authRedis,
		Jwt:             jwt,
	}
}

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

	token, err := s.Jwt.GenerateTokenAccess(u.GetUser())

	if err != nil {
		return &auth.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  "server error",
		}, nil
	}

	return &auth.LoginResponse{
		Status: http.StatusOK,
		Token:  token,
	}, nil
}

func (s Service) Validate(ctx context.Context, request *auth.ValidateRequest) (*auth.ValidateResponse, error) {
	token := request.GetToken()
	if strings.TrimSpace(token) == "" {
		return &auth.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "token not found",
		}, nil
	}

	tok, err := s.Jwt.ValidateToken(token)
	if err != nil {
		return &auth.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "token not found",
		}, nil
	}
	//check user is exists
	u, err := s.userService.GetUser(ctx, &user.UserGetRequest{
		UserId: tok.Id,
	})
	if u.Status != 200 {
		return &auth.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "user not found",
		}, nil
	}

	return &auth.ValidateResponse{
		Status: http.StatusOK,
		UserId: int64(tok.Id),
	}, nil
}

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
