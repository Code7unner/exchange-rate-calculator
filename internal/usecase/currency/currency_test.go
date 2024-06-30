package currency

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/code7unner/exchange-rate-calculator/internal/entity"
	"github.com/code7unner/exchange-rate-calculator/internal/usecase/currency/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestUseCase_Update(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		pair *entity.CurrencyPair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:  context.Background(),
				pair: entity.NewCurrencyPair(gofakeit.CurrencyShort(), gofakeit.CurrencyShort()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goMockCtrl := gomock.NewController(t)
			currencyRepoRepository := mock.NewMockcurrencyRepo(goMockCtrl)
			fastForexRepository := mock.NewMockfastForexRepo(goMockCtrl)

			rate := gofakeit.Float64()

			fastForexRepository.EXPECT().
				FetchOne(gomock.Any(), tt.args.pair.From, tt.args.pair.To).
				Return(rate, nil)

			currencyRepoRepository.EXPECT().
				UpdateExchangeRate(gomock.Any(), entity.NewExchangeRate(tt.args.pair.From, tt.args.pair.To, rate)).
				Return(nil)

			c := NewUseCase(currencyRepoRepository, fastForexRepository)
			if err := c.Update(tt.args.ctx, tt.args.pair); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCase_Convert(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		ea  *entity.ExchangeAmount
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				ea: entity.NewExchangeAmount(
					gofakeit.CurrencyShort(),
					gofakeit.CurrencyShort(),
					2.0),
			},
			want:    6.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goMockCtrl := gomock.NewController(t)
			currencyRepoRepository := mock.NewMockcurrencyRepo(goMockCtrl)

			currencyRepoRepository.EXPECT().
				GetExchangeRate(gomock.Any(), tt.args.ea.From, tt.args.ea.To).
				Return(3.0, nil)

			c := NewUseCase(currencyRepoRepository, nil)
			got, err := c.Convert(tt.args.ctx, tt.args.ea)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCase_Init(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		pair *entity.CurrencyPair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx:  context.Background(),
				pair: entity.NewCurrencyPair(gofakeit.CurrencyShort(), gofakeit.CurrencyShort()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goMockCtrl := gomock.NewController(t)
			currencyRepoRepository := mock.NewMockcurrencyRepo(goMockCtrl)
			fastForexRepository := mock.NewMockfastForexRepo(goMockCtrl)

			rate := gofakeit.Float64()

			fastForexRepository.EXPECT().
				FetchOne(gomock.Any(), tt.args.pair.From, tt.args.pair.To).
				Return(rate, nil)

			currencyRepoRepository.EXPECT().
				Upsert(gomock.Any(), entity.NewExchangeRate(tt.args.pair.From, tt.args.pair.To, rate)).
				Return(nil)

			c := NewUseCase(currencyRepoRepository, fastForexRepository)
			if err := c.Init(tt.args.ctx, tt.args.pair); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
