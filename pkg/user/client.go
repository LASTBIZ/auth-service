package user

import (
	"context"
	"google.golang.org/grpc"
	"lastbiz/auth-service/internal/config"
	"lastbiz/auth-service/pkg/logging"
)

func InitServiceClient(ctx context.Context, cfg *config.Config) UserServiceClient {
	cc, err := grpc.Dial(cfg.UserServiceURL, grpc.WithInsecure())
	if err != nil {
		logging.Error(ctx, "Could not correct:", err)
		panic(err)
	}
	return NewUserServiceClient(cc)
}
