package currency

import (
	"context"
	"github.com/code7unner/exchange-rate-calculator/internal/entity"
)

type UseCase interface {
	Convert(ctx context.Context, ea *entity.ExchangeAmount) (float64, error)
}

type Controller struct {
	useCase UseCase
}

func NewController(useCase UseCase) *Controller {
	return &Controller{useCase: useCase}
}
