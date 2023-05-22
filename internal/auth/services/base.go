package services

import (
	auth1 "lastbiz/auth-service/internal/auth"
	"lastbiz/auth-service/internal/password"
	"lastbiz/auth-service/internal/provider"
	"lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/project"
	"lastbiz/auth-service/pkg/user"
	"time"
)

type Service struct {
	passService          password.Service
	providerService      provider.Service
	PrivateKeyAccess     string
	PrivateKeyRefresh    string
	userService          user.UserServiceClient
	authRedis            auth1.Redis
	projectService       project.ProjectServiceClient
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	auth.UnimplementedAuthServiceServer
}

func NewAuthService(
	passService password.Service,
	providerService provider.Service,
	userService user.UserServiceClient,
	authRedis auth1.Redis,
	PrivateKeyAccess string,
	PrivateKeyRefresh string,
	projectService project.ProjectServiceClient,
	AccessTokenDuration time.Duration,
	RefreshTokenDuration time.Duration,
) auth.AuthServiceServer {
	return Service{
		passService:       passService,
		providerService:   providerService,
		userService:       userService,
		authRedis:         authRedis,
		PrivateKeyRefresh: PrivateKeyRefresh,
		PrivateKeyAccess:  PrivateKeyAccess,
		projectService:    projectService,
	}
}
