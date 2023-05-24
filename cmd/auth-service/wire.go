//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"auth-service/internal/biz"
	"auth-service/internal/conf"
	"auth-service/internal/data"
	"auth-service/internal/provider"
	"auth-service/internal/server"
	"auth-service/internal/service"
	"auth-service/internal/token"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Auth, *conf.Providers, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, token.NewJwtClaims, provider.Set, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
