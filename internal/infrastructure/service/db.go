package service

import (
	"context"
	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/code7unner/exchange-rate-calculator/internal/infrastructure/config"
	"github.com/code7unner/exchange-rate-calculator/third_party/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PostgresDB struct {
	logger    *zerolog.Logger
	pool      *pgxpool.Pool
	poolCheck gosundheit.Check
}

func NewPostgresDB(
	ctx context.Context,
	cfg config.PostgresConfig,
	logger *zerolog.Logger,
	startupCheckSvc *StartupCheckService,
) *PostgresDB {
	pdb := &PostgresDB{logger: logger}

	pool, err := pgxpool.New(ctx, cfg.GetURL())
	if err != nil {
		logger.Fatal().Err(err).Msgf("%s initializing error", utils.GetCurrentFuncName())
	}

	pdb.poolCheck = NewPostgresCheck(utils.GetTypeNameByObject(pdb), pool)
	err = startupCheckSvc.GetChecker().RegisterCheck(pdb.poolCheck, gosundheit.ExecutionPeriod(startupCheckSvc.GetCheckPeriod()))
	if err != nil {
		logger.Fatal().Err(err).Msgf("%s initializing error", utils.GetCurrentFuncName())
	}

	pdb.pool = pool

	logger.Info().Msgf("%s initialized successfully", utils.GetCurrentFuncName())

	return pdb
}

func (d PostgresDB) GetPool() *pgxpool.Pool {
	return d.pool
}
