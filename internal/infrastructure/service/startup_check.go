package service

import (
	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/code7unner/exchange-rate-calculator/internal/infrastructure/config"
	"time"
)

type StartupCheckService struct {
	startupChecker gosundheit.Health
	checkPeriod    time.Duration
}

func NewStartupCheck(startupCheckCfg config.CheckConfig) *StartupCheckService {
	return &StartupCheckService{
		startupChecker: gosundheit.New(),
		checkPeriod:    time.Duration(startupCheckCfg.PeriodSeconds) * time.Second,
	}
}

func (s *StartupCheckService) GetCheckPeriod() time.Duration {
	return s.checkPeriod
}

func (s *StartupCheckService) GetChecker() gosundheit.Health {
	return s.startupChecker
}
