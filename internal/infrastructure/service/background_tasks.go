package service

import (
	"context"
	"fmt"
	"github.com/code7unner/exchange-rate-calculator/internal/entity"
	"github.com/code7unner/exchange-rate-calculator/internal/infrastructure/config"
	"github.com/code7unner/exchange-rate-calculator/third_party/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"time"
)

type UseCase interface {
	Update(ctx context.Context, er *entity.CurrencyPair) error
	Init(ctx context.Context, pair *entity.CurrencyPair) error
}

type BackgroundTasks struct {
	pairs    []*entity.CurrencyPair
	tickTime time.Duration
	logger   *zerolog.Logger
	useCase  UseCase
}

func NewBackgroundTasks(cfg config.BackgroundTasksConfig, logger *zerolog.Logger, useCase UseCase) *BackgroundTasks {
	currencyPairs, err := cfg.GetCurrencyPairs()
	if err != nil {
		logger.Fatal().Err(err).Msg(utils.GetFuncName(cfg.GetCurrencyPairs))
	}

	pairs := make([]*entity.CurrencyPair, 0)
	for from, cp := range currencyPairs {
		for _, to := range cp {
			pairs = append(pairs, entity.NewCurrencyPair(from, to))
		}
	}

	logger.Info().Msgf("%s initialized successfully", utils.GetCurrentFuncName())

	return &BackgroundTasks{
		logger:   logger,
		pairs:    pairs,
		useCase:  useCase,
		tickTime: time.Duration(cfg.PeriodMinutes) * time.Minute,
	}
}

func (b *BackgroundTasks) Start(ctx context.Context) {
	// init pairs before tickers start
	if err := b.init(ctx); err != nil {
		b.logger.Fatal().Err(err).Msg(utils.GetFuncName(b.init))
	}
	for _, pair := range b.pairs {
		go b.start(ctx, pair)
	}
}

func (b *BackgroundTasks) init(ctx context.Context) error {
	for _, pair := range b.pairs {
		if err := b.useCase.Init(ctx, pair); err != nil {
			return errors.Wrap(err, fmt.Sprintf("%+v pair init failed", pair))
		}
	}
	return nil
}

func (b *BackgroundTasks) start(ctx context.Context, pair *entity.CurrencyPair) {
	ticker := time.NewTicker(b.tickTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := b.useCase.Update(ctx, pair); err != nil {
				b.logger.Error().Err(err).Msg(utils.GetFuncName(b.useCase.Update))
			}
		case <-ctx.Done():
			return
		}
	}
}
