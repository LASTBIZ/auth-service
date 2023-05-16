package auth

import (
	"context"
	"github.com/redis/go-redis/v9"
	"lastbiz/auth-service/internal/utils"
)

type authRedis struct {
	redis *redis.Client
}

func NewAuthRedis(redis *redis.Client) *authRedis {
	return &authRedis{
		redis: redis,
	}
}

func (a authRedis) CreateState(ctx context.Context) (string, error) {
	state := utils.GenerateState()
	//SET key value EX seconds
	_, err := a.redis.Set(ctx, "auth:state:"+state, state, 600).Result()
	if err != nil {
		return "", err
	}
	return state, nil
}

func (a authRedis) CheckState(ctx context.Context, state string) (bool, error) {
	result, err := a.redis.Exists(ctx, "auth:state"+state).Result()
	if err != nil {
		return false, err
	}
	if result >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}
