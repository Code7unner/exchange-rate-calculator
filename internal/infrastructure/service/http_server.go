package service

import (
	"fmt"
	"github.com/code7unner/exchange-rate-calculator/internal/adapter/controller/http/currency"
	"github.com/code7unner/exchange-rate-calculator/internal/infrastructure/config"
	"github.com/code7unner/exchange-rate-calculator/third_party/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/rs/zerolog"
)

type HTTPServer struct {
	port int
	app  *fiber.App
}

func NewHTTPServer(
	httpCfg config.HTTPServerConfig,
	logger *zerolog.Logger,
	startupCheckSvc *StartupCheckService,
	currencyController *currency.Controller,
) *HTTPServer {
	app := fiber.New()
	app.Server().Logger = logger

	app.Get("/startupz", utils.HandleHealthJSON(startupCheckSvc.GetChecker()))
	v1 := app.Group("/v1")
	v1.Use(cors.New())
	v1.Get("/currency", currencyController.Convert)

	logger.Info().Msgf("%s initialized successfully", utils.GetCurrentFuncName())

	return &HTTPServer{app: app, port: httpCfg.Port}
}

func (s *HTTPServer) Start() error {
	return s.app.Listen(fmt.Sprintf(":%d", s.port))
}
