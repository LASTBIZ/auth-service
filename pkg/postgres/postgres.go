package postgres

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lastbiz/auth-service/internal/config"
	"lastbiz/auth-service/pkg/utils"
	"log"
	"time"
)

func NewClient(ctx context.Context, maxAttemps int, sc config.Postgres) *gorm.DB {
	var pool *gorm.DB
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", sc.Host, sc.User, sc.Password, sc.DB, sc.Port)
	err = utils.DoWithTries(func() error {
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		return nil
	}, maxAttemps, 5*time.Second)

	if err != nil {
		log.Fatal("error to connect in PostgreSQL max attemtps")
	}
	return pool
}
