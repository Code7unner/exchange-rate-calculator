package currency

import (
	"context"
	"github.com/code7unner/exchange-rate-calculator/internal/entity"
	"github.com/code7unner/exchange-rate-calculator/third_party/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	PairNotFoundError = errors.New("currency pair not found")
)

type PostgresRepository struct {
	repo   *pgxpool.Pool
	logger *zerolog.Logger
}

func NewPostgresRepository(
	repo *pgxpool.Pool,
	logger *zerolog.Logger,
) *PostgresRepository {
	return &PostgresRepository{
		repo:   repo,
		logger: logger,
	}
}

func (r *PostgresRepository) UpdateExchangeRate(ctx context.Context, er *entity.ExchangeRate) error {
	const query = `
		UPDATE exchange_rates
		SET rate = $1, updated_at = CURRENT_TIMESTAMP
		WHERE from_currency = $2 AND to_currency = $3;`

	tx, err := r.repo.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, utils.GetFuncName(r.repo.Begin))
	}

	_, err = tx.Exec(ctx, query, er.Rate, er.From, er.To)
	if err != nil {
		return errors.Wrap(err, utils.GetFuncName(tx.Exec))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, utils.GetFuncName(tx.Commit))
	}

	return nil
}

func (r *PostgresRepository) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	const query = `
		SELECT rate 
		FROM exchange_rates 
		WHERE from_currency = $1 AND to_currency = $2;`

	var rate float64
	if err := r.repo.QueryRow(ctx, query, from, to).Scan(&rate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, PairNotFoundError
		}
		return 0, errors.Wrap(err, utils.GetFuncName(r.repo.QueryRow))
	}

	return rate, nil
}

func (r *PostgresRepository) Upsert(ctx context.Context, er *entity.ExchangeRate) error {
	const query = `INSERT INTO exchange_rates (from_currency, to_currency, rate, created_at, updated_at) 
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (from_currency, to_currency) 
		DO UPDATE SET rate = EXCLUDED.rate, updated_at = CURRENT_TIMESTAMP;`

	tx, err := r.repo.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, utils.GetFuncName(r.repo.Begin))
	}

	_, err = tx.Exec(ctx, query, er.From, er.To, er.Rate)
	if err != nil {
		return errors.Wrap(err, utils.GetFuncName(tx.Exec))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, utils.GetFuncName(tx.Commit))
	}

	return nil
}
