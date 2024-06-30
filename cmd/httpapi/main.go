package main

import (
	"context"
	"github.com/code7unner/exchange-rate-calculator/internal/infrastructure/app"
	"github.com/code7unner/exchange-rate-calculator/internal/infrastructure/config"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	cfg := config.Get(logger)
	app.Start(ctx, cfg, logger)
}
