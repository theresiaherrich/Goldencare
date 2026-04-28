package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/bootstrap"
	"github.com/theresiaherrich/Goldencare/internal/middleware"
)

func registerSuperadminRoutes(api fiber.Router, deps *AppDependencies, container *bootstrap.Container) {
	superadmin := api.Group("/superadmin",
		middleware.RequireAuth(container.Config),
		middleware.RequireRole("superadmin"),
	)
	superadmin.Post("/kode/generate", deps.AuthHandler.SuperadminGenerateKode)
	superadmin.Get("/panti", deps.PantiHandler.GetAll)
	superadmin.Post("/panti", deps.PantiHandler.Create)
}