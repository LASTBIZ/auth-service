package config

import (
	"lastbiz/auth-service/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort       string `yaml:"grpc_port" env:"GRPC_PORT" env-required:"true"`
	UserServiceURL string `yaml:"userServiceURL" env:"USER_SERVICE_URL" env-required:"true"`
	JWT            struct {
		SecretKeyAccess        string `yaml:"secret_key_access" env:"SECRET_KEY_ACCESS" env-required:"true"`
		SecretKeyRefresh       string `yaml:"secret_key_refresh" env:"SECRET_KEY_REFRESH" env-required:"true"`
		ExpirationHoursAccess  string `yaml:"expiry_access_token" env:"EXPIRY_ACCESS_TOKEN" env-required:"true"`
		ExpirationHoursRefresh string `yaml:"expiry_refresh_token" env:"EXPIRY_REFRESH_TOKEN" env-required:"true"`
	} `yaml:"jwt"`
	Postgres struct {
		Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
		User     string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
		Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
		DB       string `yaml:"db" env:"POSTGRES_DATABASE" env-required:"true"`
		Port     string `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	} `yaml:"postgresql"`
	Redis struct {
		Host     string `yaml:"host" env:"REDIS_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"REDIS_PORT" env-required:"true"`
		Password string `yaml:"password" env:"REDIS_PASSWORD" env-required:"true"`
	} `yaml:"redis"`
	Providers struct {
		Google struct {
			ClientID         string `yaml:"client_id" env:"GOOGLE_CLIENT_ID" env-required:"true"`
			ClientSecret     string `yaml:"client_secret" env:"GOOGLE_CLIENT_SECRET" env-required:"true"`
			OAuthRedirectURl string `yaml:"o_auth_redirect_u_rl" env:"GOOGLE_OAUTH_REDIRECT_URL" env-required:"true"`
		} `yaml:"google"`
		Facebook struct {
			ClientID         string `yaml:"client_id" env:"FACEBOOK_CLIENT_ID" env-required:"true"`
			ClientSecret     string `yaml:"client_secret" env:"FACEBOOK_CLIENT_SECRET" env-required:"true"`
			OAuthRedirectURl string `yaml:"o_auth_redirect_u_rl" env:"FACEBOOK_OAUTH_REDIRECT_URL" env-required:"true"`
		} `yaml:"facebook"`
	} `yaml:"providers"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			helpText := "LastMBiz auth-service by https://github.com/Suro4ek"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
