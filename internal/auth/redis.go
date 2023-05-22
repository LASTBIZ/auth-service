package auth

import (
	"context"
	"github.com/redis/go-redis/v9"
	"lastbiz/auth-service/internal/utils"
)

type Redis struct {
	redis *redis.Client
}

func NewAuthRedis(redis *redis.Client) *Redis {
	return &Redis{
		redis: redis,
	}
}

func (r Redis) CreateState(ctx context.Context) (string, error) {
	state := utils.GenerateState()
	//SET key value EX seconds
	_, err := r.redis.Set(ctx, "auth:state:"+state, state, 600).Result()
	if err != nil {
		return "", err
	}
	return state, nil
}

func (r Redis) CheckState(ctx context.Context, state string) (bool, error) {
	result, err := r.redis.Exists(ctx, "auth:state"+state).Result()
	if err != nil {
		return false, err
	}
	if result >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}
