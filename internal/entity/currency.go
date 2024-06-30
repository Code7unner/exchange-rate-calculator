package entity

type CurrencyPair struct {
	From string
	To   string
}

func NewCurrencyPair(from string, to string) *CurrencyPair {
	return &CurrencyPair{From: from, To: to}
}

type ExchangeRate struct {
	From string
	To   string
	Rate float64
}

func NewExchangeRate(from string, to string, rate float64) *ExchangeRate {
	return &ExchangeRate{From: from, To: to, Rate: rate}
}

type ExchangeAmount struct {
	From   string
	To     string
	Amount float64
}

func NewExchangeAmount(from string, to string, amount float64) *ExchangeAmount {
	return &ExchangeAmount{From: from, To: to, Amount: amount}
}
