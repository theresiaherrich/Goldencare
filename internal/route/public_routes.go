package route

import "github.com/gofiber/fiber/v2"

func registerPublicRoutes(api fiber.Router, deps *AppDependencies) {
	api.Post("/auth/register", deps.AuthHandler.Register)
	api.Post("/auth/login", deps.AuthHandler.Login)
	api.Post("/superadmin/login", deps.AuthHandler.SuperadminLogin)
}