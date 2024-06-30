package currency

import (
	"context"
	"github.com/code7unner/exchange-rate-calculator/internal/entity"
	"github.com/code7unner/exchange-rate-calculator/third_party/utils"
	"github.com/pkg/errors"
)

//go:generate mockgen -source=currency.go -destination=mock/currency.go -package=mock

type currencyRepo interface {
	UpdateExchangeRate(ctx context.Context, er *entity.ExchangeRate) error
	GetExchangeRate(ctx context.Context, from, to string) (float64, error)
	Upsert(ctx context.Context, er *entity.ExchangeRate) error
}

type fastForexRepo interface {
	FetchOne(ctx context.Context, from, to string) (float64, error)
}

type UseCase struct {
	currencyRepo  currencyRepo
	fastForexRepo fastForexRepo
}

func NewUseCase(currencyRepo currencyRepo, fastForexRepo fastForexRepo) *UseCase {
	return &UseCase{currencyRepo: currencyRepo, fastForexRepo: fastForexRepo}
}

func (c *UseCase) Update(ctx context.Context, pair *entity.CurrencyPair) error {
	rate, err := c.fastForexRepo.FetchOne(ctx, pair.From, pair.To)
	if err != nil {
		return errors.Wrap(err, utils.GetFuncName(c.fastForexRepo.FetchOne))
	}

	return c.currencyRepo.UpdateExchangeRate(
		ctx,
		entity.NewExchangeRate(pair.From, pair.To, rate),
	)
}

func (c *UseCase) Convert(ctx context.Context, ea *entity.ExchangeAmount) (float64, error) {
	rate, err := c.currencyRepo.GetExchangeRate(ctx, ea.From, ea.To)
	if err != nil {
		return 0, errors.Wrap(err, utils.GetFuncName(c.currencyRepo.GetExchangeRate))
	}

	totalAmount := rate * ea.Amount

	return totalAmount, nil
}

func (c *UseCase) Init(ctx context.Context, pair *entity.CurrencyPair) error {
	rate, err := c.fastForexRepo.FetchOne(ctx, pair.From, pair.To)
	if err != nil {
		return errors.Wrap(err, utils.GetFuncName(c.fastForexRepo.FetchOne))
	}

	return c.currencyRepo.Upsert(
		ctx,
		entity.NewExchangeRate(pair.From, pair.To, rate),
	)
}
