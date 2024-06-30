package service

import (
	"context"
	"github.com/code7unner/exchange-rate-calculator/third_party/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PostgresCheck struct {
	checkName string
	pool      *pgxpool.Pool
}

func NewPostgresCheck(name string, pool *pgxpool.Pool) *PostgresCheck {
	return &PostgresCheck{
		checkName: name,
		pool:      pool,
	}
}

func (check *PostgresCheck) Name() string {
	return check.checkName
}

func (check *PostgresCheck) Execute(ctx context.Context) (interface{}, error) {
	if err := check.pool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, utils.GetFuncName(check.pool.Ping))
	}

	return nil, nil
}
