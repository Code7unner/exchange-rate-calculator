package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/code7unner/exchange-rate-calculator/third_party/utils"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const appName = "exchange-rate-calculator"

type CheckConfig struct {
	PeriodSeconds int `envconfig:"PERIOD_SECONDS" required:"true"`
}

type HTTPServerConfig struct {
	Port int `envconfig:"PORT" required:"true"`
}

type HTTPClientConfig struct {
	Host     string `envconfig:"HOST" required:"true"`
	Endpoint string `envconfig:"ENDPOINT" required:"true"`
}

type BackgroundTasksConfig struct {
	PeriodMinutes          int `envconfig:"PERIOD_MINUTES" required:"true"`
	SupportedCurrencyPairs map[string]string
}

//go:embed currency-pairs.json
var currencyPairsFile []byte

func (b BackgroundTasksConfig) GetCurrencyPairs() (map[string][]string, error) {
	var currencyPairs map[string][]string
	if err := json.Unmarshal(currencyPairsFile, &currencyPairs); err != nil {
		return nil, errors.Wrap(err, utils.GetFuncName(json.Unmarshal))
	}
	return currencyPairs, nil
}

type PostgresConfig struct {
	DB       string `envconfig:"DB" required:"true"`
	Host     string `envconfig:"HOST" required:"true"`
	User     string `envconfig:"USER" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Port     int    `envconfig:"PORT" required:"true"`
}

func (p PostgresConfig) GetURL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", p.User, p.Password, p.Host, p.Port, p.DB)
}

func (p PostgresConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", p.Host, p.Port)
}

type FastForexConfig struct {
	APIKey string `envconfig:"API_KEY" required:"true"`
}

type Config struct {
	AppName              string
	LogLevel             string                `envconfig:"LOG_LEVEL" required:"true"`
	ShutdownPauseSeconds int                   `envconfig:"SHUTDOWN_PAUSE_SECONDS" required:"true"`
	StartupCheck         CheckConfig           `envconfig:"STARTUP_CHECK" required:"true"`
	Postgres             PostgresConfig        `envconfig:"POSTGRES" required:"true"`
	FastForex            FastForexConfig       `envconfig:"FAST_FOREX" required:"true"`
	HTTPServer           HTTPServerConfig      `envconfig:"HTTP_SERVER" required:"true"`
	BackgroundTasks      BackgroundTasksConfig `envconfig:"BACKGROUND_TASKS" required:"true"`
}

func Get(logger zerolog.Logger) *Config {
	cfg := &Config{AppName: appName}
	if err := envconfig.Process("", cfg); err != nil {
		logger.Fatal().Err(err).Msg(utils.GetFuncName(envconfig.Process))
	}
	return cfg
}
