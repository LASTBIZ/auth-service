package biz

import (
	"context"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewHashUseCase, NewProviderUseCase, NewAuthUsecase)

//go:generate mockgen -destination=../mocks/usecase/biz.go -package=usecase . Transaction
type Transaction interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}
