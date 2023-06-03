package biz

import (
	"auth-service/api/investor"
	"auth-service/api/user"
	"auth-service/internal/provider"
	"auth-service/internal/token"
	"auth-service/internal/utils"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// AuthUsecase is a Auth usecase.
type AuthUseCase struct {
	uh       *HashUseCase
	up       *ProviderUseCase
	uc       user.UserClient
	ic       investor.InvestorClient
	log      *log.Helper
	provider *provider.Struct
	claims   *token.JwtClaims
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

// NewAuthUsecase new a Auth usecase.
func NewAuthUsecase(uh *HashUseCase, uc user.UserClient, up *ProviderUseCase, ic investor.InvestorClient, provider *provider.Struct, claims *token.JwtClaims, logger log.Logger) *AuthUseCase {
	return &AuthUseCase{uh: uh, uc: uc, ic: ic, log: log.NewHelper(logger), up: up, provider: provider, claims: claims}
}

func (au *AuthUseCase) Register(ctx context.Context, email, firstName, lastName, password string) (bool, error) {
	u, err := au.uc.CreateUser(ctx, &user.CreateUserRequest{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	})

	if err != nil {
		return false, err
	}

	_, err = au.ic.CreateInvestor(ctx, &investor.CreateInvestorRequest{
		FullName: u.FirstName + " " + u.LastName,
		UserId:   uint64(u.Id),
	})

	if err != nil {
		return false, err
	}

	password = utils.HashPassword(password)
	_, err = au.uh.Create(ctx, &Hash{
		UserID: u.Id,
		Hash:   password,
	})

	if err != nil {
		_, err := au.uc.DeleteUser(ctx, &user.IdRequest{
			Id: int64(u.Id),
		})
		_, err = au.ic.DeleteInvestor(ctx, &investor.DeleteInvestorRequest{
			Id: uint64(u.Id),
		})
		return false, err
	}

	return true, nil
}

func (au *AuthUseCase) CreateState(provider string) (string, error) {
	prov, ok := au.provider.Providers[provider]
	if !ok {
		return "", errors.NotFound("PROVIDER_NOT_FOUND", "provider not found")
	}

	state, err := au.up.CreateState(context.Background())

	if err != nil {
		return "", err
	}

	return prov.GenerateOAuthToken(state), nil
}

func (au *AuthUseCase) Login(ctx context.Context, email, password string) (*Token, error) {
	u, err := au.uc.GetUserByEmail(ctx, &user.EmailRequest{
		Email: email,
	})

	if err != nil {
		return nil, err
	}

	hash, err := au.uh.GetHashByUserId(ctx, u.Id)

	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(password, hash.Hash) {
		return nil, errors.Unauthorized("WRONG_PASSWORD", "password is not match")
	}

	access, err := au.claims.Access.CreateToken(hash.UserID)
	if err != nil {
		return nil, errors.InternalServer("ERROR_CREATE_TOKEN", "error create token")
	}

	refresh, err := au.claims.Refresh.CreateToken(hash.UserID)
	if err != nil {
		return nil, errors.InternalServer("ERROR_CREATE_TOKEN", "error create token")
	}

	token := &Token{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	return token, nil
}

func (au *AuthUseCase) Callback(ctx context.Context, provider, code, state string) (*Token, error) {
	prov, ok := au.provider.Providers[provider]
	if !ok {
		return nil, errors.NotFound("PROVIDER_NOT_FOUND", "provider not found")
	}

	err := au.up.CheckState(ctx, state)
	if err != nil {
		return nil, errors.NotFound("STATE_NOT_FOUND", "state not found")
	}

	token, err := prov.Callback(code)
	if err != nil {
		return nil, errors.New(505, "PROVIDER_CALLBACK_ERROR", "provider callback error")
	}

	userProv, err := prov.GetUser(token)
	if err != nil {
		return nil, errors.New(505, "PROVIDER_CALLBACK_ERROR", "provider callback error")
	}

	//create user
	_, err = au.up.GetProviderByEmail(ctx, userProv.Email)
	if err != nil {
		if errors.IsNotFound(err) {
			//register user
			u, err := au.uc.CreateUser(ctx, &user.CreateUserRequest{
				Email:     userProv.Email,
				FirstName: userProv.GivenName,
				LastName:  userProv.FamilyName,
			})

			if err != nil {
				return nil, err
			}

			_, err = au.up.Create(ctx, &Provider{
				Email:        userProv.Email,
				UserID:       u.Id,
				Provider:     provider,
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				TokenType:    token.TokenType,
			})

			if err != nil {
				//TODO delete user
				return nil, err
			}
		}
		return nil, err
	}

	u, err := au.uc.GetUserByEmail(ctx, &user.EmailRequest{
		Email: userProv.Email,
	})

	_, err = au.uc.UpdateUser(ctx, &user.UpdateUserRequest{
		Id:           int64(u.Id),
		FirstName:    userProv.GivenName,
		LastName:     userProv.FamilyName,
		Role:         u.Role,
		IsVerify:     userProv.VerifyEmail,
		Phone:        u.Phone,
		Organization: u.Organization,
		Messenger:    u.Messengers,
		Blocked:      u.Blocked,
	})

	_, err = au.up.Update(ctx, &Provider{
		Email:        userProv.Email,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	})

	if err != nil {
		return nil, err
	}

	access, err := au.claims.Access.CreateToken(u.Id)
	if err != nil {
		return nil, errors.InternalServer("ERROR_CREATE_TOKEN", "error create token")
	}

	refresh, err := au.claims.Refresh.CreateToken(u.Id)
	if err != nil {
		return nil, errors.InternalServer("ERROR_CREATE_TOKEN", "error create token")
	}

	tk := &Token{
		AccessToken:  access,
		RefreshToken: refresh,
	}
	return tk, nil
}

func (au *AuthUseCase) Validate(ctx context.Context, token string) (float64, error) {
	id, err := au.claims.Access.ValidateToken(token)
	if err != nil {
		return 0, errors.Unauthorized("WRONG_TOKEN", "wrong token")
	}
	return id.(float64), nil
}

func (au *AuthUseCase) RefreshToken(refreshToken string) (*Token, error) {
	id, err := au.claims.Refresh.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.Unauthorized("WRONG_TOKEN", "wrong token")
	}

	access, err := au.claims.Access.CreateToken(id)
	if err != nil {
		return nil, errors.InternalServer("ERROR_CREATE_TOKEN", "error create token")
	}

	refresh, err := au.claims.Refresh.CreateToken(id)
	if err != nil {
		return nil, errors.InternalServer("ERROR_CREATE_TOKEN", "error create token")
	}

	tk := &Token{
		AccessToken:  access,
		RefreshToken: refresh,
	}
	return tk, nil
}

func (au *AuthUseCase) ChangePassword(ctx context.Context, userId uint32, password, newPassword string) (bool, error) {
	u, err := au.uc.GetUserById(ctx, &user.IdRequest{
		Id: int64(userId),
	})

	if err != nil {
		return false, err
	}

	hash, err := au.uh.GetHashByUserId(ctx, u.Id)

	if err != nil {
		if errors.IsNotFound(err) && password == "" {
			_, err := au.uh.Create(ctx, &Hash{
				UserID: u.Id,
				Hash:   utils.HashPassword(newPassword),
			})
			if err != nil {
				return false, err
			}
		}
		return false, err
	}

	if !utils.CheckPasswordHash(password, hash.Hash) {
		return false, errors.Unauthorized("WRONG_PASSWORD", "password is not match")
	}

	_, err = au.uh.Update(ctx, &Hash{ID: hash.ID, Hash: utils.HashPassword(newPassword)})

	if err != nil {
		return false, err
	}

	return true, nil
}
