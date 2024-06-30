package app

import (
	"context"
	"github.com/code7unner/exchange-rate-calculator/internal/adapter/controller/http/currency"
	currencyRepo "github.com/code7unner/exchange-rate-calculator/internal/adapter/repository/currency"
	"github.com/code7unner/exchange-rate-calculator/internal/adapter/repository/fast-forex"
	"github.com/code7unner/exchange-rate-calculator/internal/infrastructure/config"
	"github.com/code7unner/exchange-rate-calculator/internal/infrastructure/service"
	currencyUC "github.com/code7unner/exchange-rate-calculator/internal/usecase/currency"
	"github.com/code7unner/exchange-rate-calculator/third_party/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func Start(ctx context.Context, cfg *config.Config, logger zerolog.Logger) {
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.Fatal().
			Err(errors.Wrap(err, utils.GetFuncName(zerolog.ParseLevel)))
	}
	logger.Level(logLevel)

	logger.Info().Msg("Application starting...")

	startupCheckSvc := service.NewStartupCheck(cfg.StartupCheck)
	db := service.NewPostgresDB(ctx, cfg.Postgres, &logger, startupCheckSvc)

	currencyRepository := currencyRepo.NewPostgresRepository(db.GetPool(), &logger)
	fastForexRepository := fastforex.NewHTTPRepository(cfg.FastForex.APIKey)

	currencyUseCase := currencyUC.NewUseCase(currencyRepository, fastForexRepository)

	currencyCtrl := currency.NewController(currencyUseCase)

	httpSvc := service.NewHTTPServer(cfg.HTTPServer, &logger, startupCheckSvc, currencyCtrl)

	backgroundTasks := service.NewBackgroundTasks(cfg.BackgroundTasks, &logger, currencyUseCase)

	// Start background tasks
	backgroundTasks.Start(ctx)

	// Start HTTP server
	logger.Fatal().
		Err(httpSvc.Start())
}
