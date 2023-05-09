package config

import (
	"lastbiz/auth-service/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort       string `yaml:"grpc_port" env:"GRPC_PORT"`
	UserServiceURL string `yaml:"userServiceURL" env:"USER_SERVICE_URL"`
}

type Postgres struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	User     string `yaml:"host" env:"POSTGRES_USER"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
	DB       string `yaml:"db" env:"POSTGRES_DATABASE"`
	Port     string `yaml:"port" env:"POSTGRES_PORT"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
