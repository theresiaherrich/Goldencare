package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/bootstrap"
)

func registerHealthRoutes(api fiber.Router, container *bootstrap.Container) {
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "OK",
			"app":    container.Config.AppName,
		})
	})
}