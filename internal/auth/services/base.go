package services

import (
	auth1 "lastbiz/auth-service/internal/auth"
	"lastbiz/auth-service/internal/password"
	"lastbiz/auth-service/internal/provider"
	"lastbiz/auth-service/internal/utils"
	"lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/project"
	"lastbiz/auth-service/pkg/user"
)

type Service struct {
	passService     password.Service
	providerService provider.Service
	Jwt             utils.JwtWrapper
	userService     user.UserServiceClient
	authRedis       auth1.Redis
	projectService  project.ProjectServiceClient
	auth.UnimplementedAuthServiceServer
}

func NewAuthService(
	passService password.Service,
	providerService provider.Service,
	userService user.UserServiceClient,
	authRedis auth1.Redis,
	jwt utils.JwtWrapper,
	projectService project.ProjectServiceClient,
) auth.AuthServiceServer {
	return Service{
		passService:     passService,
		providerService: providerService,
		userService:     userService,
		authRedis:       authRedis,
		Jwt:             jwt,
		projectService:  projectService,
	}
}
