package main

import (
	"context"
	"google.golang.org/grpc"
	"lastbiz/auth-service/internal/auth"
	"lastbiz/auth-service/internal/config"
	"lastbiz/auth-service/internal/password"
	"lastbiz/auth-service/internal/provider"
	"lastbiz/auth-service/internal/provider/providers"
	"lastbiz/auth-service/internal/utils"
	auth2 "lastbiz/auth-service/pkg/auth"
	"lastbiz/auth-service/pkg/logging"
	"lastbiz/auth-service/pkg/postgres"
	"lastbiz/auth-service/pkg/redis"
	"lastbiz/auth-service/pkg/user"
	"net"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logging.Info(ctx, "config initializing")
	cfg := config.GetConfig()

	pgconfig := postgres.NewPGConfig(cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DB)

	pgClient := postgres.NewClient(ctx, 5, time.Second*5, pgconfig)

	redisConfig := redis.NewRedisConfig(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	redisClient, err := redis.NewClient(context.Background(), *redisConfig)
	if err != nil {
		logging.Error(ctx, "error connect redis: ", err)
	}

	passStorage := password.NewPasswordStorage(*pgClient)
	passService := password.NewPasswordService(*passStorage)

	providerStorage := provider.NewProviderStorage(*pgClient)
	prs := make(map[string]provider.Provider, 0)
	prs["google"] = providers.NewGoogleProvider(
		cfg.Providers.Google.ClientID,
		cfg.Providers.Google.ClientSecret,
		cfg.Providers.Google.OAuthRedirectURl,
		*providerStorage,
	)
	prs["facebook"] = providers.NewFacebookProvider(
		cfg.Providers.Facebook.ClientID,
		cfg.Providers.Facebook.ClientSecret,
		cfg.Providers.Facebook.OAuthRedirectURl,
		*providerStorage,
	)

	jwt := utils.JwtWrapper{
		SecretKey:              cfg.JWT.SecretKey,
		Issuer:                 cfg.JWT.Issuer,
		ExpirationHoursRefresh: int64(cfg.JWT.ExpirationHoursRefresh),
		ExpirationHoursAccess:  int64(cfg.JWT.ExpirationHoursAccess),
	}

	providerService := provider.NewProviderService(*providerStorage, prs)
	logging.Info(ctx, "connect user service")
	userService := user.InitServiceClient(ctx, cfg)
	authRedis := auth.NewAuthRedis(redisClient)
	authService := auth.NewAuthService(*passService, *providerService, userService, *authRedis, jwt)

	lis, err := net.Listen("tcp", "0.0.0.0:"+cfg.GRPCPort)

	if err != nil {
		logging.GetLogger().Fatal(err)
	}
	logging.Info(ctx, "start grpc auth server by Suro")
	grpcServer := grpc.NewServer()
	auth2.RegisterAuthServiceServer(grpcServer, authService)

	if err := grpcServer.Serve(lis); err != nil {
		logging.GetLogger().Fatal(err)
	}
}
