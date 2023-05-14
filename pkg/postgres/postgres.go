package postgres

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type pgConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func NewPGConfig(username, password, host, port, database string) *pgConfig {
	return &pgConfig{
		username,
		password,
		host,
		port,
		database,
	}
}

func NewClient(ctx context.Context, maxAttempts int, maxDelay time.Duration, cfg *pgConfig) *gorm.DB {
	var pool *gorm.DB
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port)
	err = DoWithTries(func() error {
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, maxDelay)

	if err != nil {
		log.Fatal("All attempts are exceeded. Unable to connect to postgres")
	}
	return pool
}

func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}

		return nil
	}

	return
}
