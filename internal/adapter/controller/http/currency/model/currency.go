package model

import (
	"github.com/code7unner/exchange-rate-calculator/internal/entity"
	"github.com/code7unner/exchange-rate-calculator/pkg/currency"
	"github.com/code7unner/exchange-rate-calculator/third_party/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

var (
	CryptoToCryptoConversionNotSupported = errors.New("crypto to crypto conversion is not supported")
	FiatToFiatConversionNotSupported     = errors.New("fiat to fiat conversion is not supported")
	WrongCurrencyCode                    = errors.New("wrong currency code")
	WrongCurrencyAmount                  = errors.New("wrong currency amount: must be positive")
)

type ConvertCurrencyRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func (r ConvertCurrencyRequest) Validate() error {
	if r.Amount <= 0 {
		return WrongCurrencyAmount
	}
	if r.isC2C() {
		return CryptoToCryptoConversionNotSupported
	}
	if r.isF2F() {
		return FiatToFiatConversionNotSupported
	}
	if !r.isC2F() && !r.isF2C() {
		return WrongCurrencyCode
	}

	return nil
}

func (r ConvertCurrencyRequest) ToDomain() *entity.ExchangeAmount {
	return entity.NewExchangeAmount(r.From, r.To, r.Amount)
}

func NewConvertCurrencyRequest(c fiber.Ctx) (*ConvertCurrencyRequest, error) {
	amount, err := strconv.ParseFloat(strings.TrimSpace(c.Query("amount")), 64)
	if err != nil {
		return nil, errors.Wrap(err, utils.GetFuncName(strconv.ParseFloat))
	}

	req := &ConvertCurrencyRequest{
		From:   strings.ToUpper(c.Query("from")),
		To:     strings.ToUpper(c.Query("to")),
		Amount: amount,
	}

	if err = req.Validate(); err != nil {
		return nil, errors.Wrap(err, utils.GetFuncName(req.Validate))
	}

	return req, nil
}

// isC2C: is conversion crypto-crypto
func (r ConvertCurrencyRequest) isC2C() bool {
	_, convertFromCryptoExists := currency.CryptoCurrencies[r.From]
	_, convertToCryptoExists := currency.CryptoCurrencies[r.To]
	return convertToCryptoExists && convertFromCryptoExists
}

// isF2F: is conversion fiat-fiat
func (r ConvertCurrencyRequest) isF2F() bool {
	_, convertFromFiatExists := currency.FiatCurrencies[r.From]
	_, convertToFiatExists := currency.FiatCurrencies[r.To]
	return convertFromFiatExists && convertToFiatExists
}

// isC2F: is conversion crypto-fiat
func (r ConvertCurrencyRequest) isC2F() bool {
	_, convertFromCryptoExists := currency.CryptoCurrencies[r.From]
	_, convertToFiatExists := currency.FiatCurrencies[r.To]
	return convertFromCryptoExists && convertToFiatExists
}

// isF2C: is conversion fiat-crypto
func (r ConvertCurrencyRequest) isF2C() bool {
	_, convertFromFiatExists := currency.FiatCurrencies[r.From]
	_, convertToCryptoExists := currency.CryptoCurrencies[r.To]
	return convertFromFiatExists && convertToCryptoExists
}
