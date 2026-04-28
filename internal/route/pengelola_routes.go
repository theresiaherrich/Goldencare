package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/handlers"
)

func registerPengelolaRoutes(
	r fiber.Router,
	pantiHandler *handlers.PantiHandler,
	authHandler *handlers.AuthHandler,
	lansiaHandler *handlers.LansiaHandler,
	pengurusHandler *handlers.PengurusHandler,
) {
	r.Get("/panti", pantiHandler.GetPanti)
	r.Put("/panti/:id", pantiHandler.Update)

	r.Post("/kode/generate", authHandler.GenerateKode)
	r.Get("/kode", authHandler.ListKode)
	r.Delete("/kode/:kode", authHandler.NonaktifkanKode)

	r.Get("/lansia/dashboard", lansiaHandler.GetDashboard)
	r.Get("/lansia", lansiaHandler.GetAll)
	r.Get("/lansia/:id", lansiaHandler.GetByID)
	r.Post("/lansia", lansiaHandler.Create)
	r.Put("/lansia/:id", lansiaHandler.Update)

	r.Get("/kamar", pengurusHandler.GetKamar)
	r.Post("/kamar", pengurusHandler.CreateKamar)
}