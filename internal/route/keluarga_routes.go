package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/handlers"
)

func registerKeluargaRoutes(keluarga fiber.Router, h *handlers.KeluargaHandler, lansiaHandler *handlers.LansiaHandler) {
	keluarga.Get("/dashboard/:lansia_id", h.GetDashboard)
	keluarga.Get("/ringkasan/:lansia_id", h.GetKeluarga)

	keluarga.Get("/lansia/dashboard", lansiaHandler.GetDashboard)
}