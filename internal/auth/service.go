package auth

import (
	"context"
	"lastbiz/auth-service/internal/password"
	"lastbiz/auth-service/internal/provider"
	"lastbiz/auth-service/internal/utils"
	"lastbiz/auth-service/pkg/pb/auth"
)

type Service struct {
	passService     password.Service
	providerService provider.Service
	Jwt             utils.JwtWrapper
}

func NewAuthService(passService password.Service, providerService provider.Service) *Service {
	return &Service{
		passService:     passService,
		providerService: providerService,
	}
}

func (s Service) Register(ctx context.Context, req *auth.RegisterByPasswordRequest) (*auth.RegisterResponse, error) {

}

func (s Service) RegisterByProvider(ctx context.Context) {

}

func (s Service) Login(ctx context.Context) {

}

func (s Service) LoginByProvider(ctx context.Context) {

}
