package utils

import (
	"encoding/json"
	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/gofiber/fiber/v3"
	"github.com/pkg/errors"
)

const (
	// ReportTypeShort is the value to be passed in the request parameter `type` when a short response is desired.
	ReportTypeShort = "short"
)

func HandleHealthJSON(h gosundheit.Health) fiber.Handler {
	return func(c fiber.Ctx) error {
		results, healthy := h.Results()
		c.Response().Header.Set("Content-Type", "application/json")
		if healthy {
			c.Response().SetStatusCode(200)
		} else {
			c.Response().SetStatusCode(503)
		}

		encoder := json.NewEncoder(c.Response().BodyWriter())
		encoder.SetIndent("", "\t")
		var err error

		if c.Query("type") == ReportTypeShort {
			shortResults := make(map[string]string)
			for k, v := range results {
				if v.IsHealthy() {
					shortResults[k] = "PASS"
				} else {
					shortResults[k] = "FAIL"
				}
			}

			err = encoder.Encode(shortResults)
		} else {
			err = encoder.Encode(results)
		}

		if err != nil {
			return errors.Wrap(err, GetCurrentFuncName())
		}

		return nil
	}
}
