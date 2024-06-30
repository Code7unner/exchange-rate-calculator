package currency

import (
	"errors"
	"github.com/code7unner/exchange-rate-calculator/internal/adapter/controller/http/currency/model"
	"github.com/code7unner/exchange-rate-calculator/internal/adapter/repository/currency"
	"github.com/gofiber/fiber/v3"
	"net/http"
)

func (ctrl *Controller) Convert(c fiber.Ctx) error {
	req, err := model.NewConvertCurrencyRequest(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	amount, err := ctrl.useCase.Convert(c.UserContext(), req.ToDomain())
	if err != nil {
		if errors.Is(err, currency.PairNotFoundError) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": currency.PairNotFoundError.Error()})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"currency": amount})
}
