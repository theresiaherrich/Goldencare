package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/handlers"
)

func registerKunjunganRoutes(protected fiber.Router, pengurus fiber.Router, h *handlers.KunjunganHandler) {
	protected.Get("/kunjungan/lansia/:lansia_id", h.GetByLansia)
	protected.Get("/kunjungan/lansia/:lansia_id/terbaru", h.GetTerbaru)

	pengurus.Post("/kunjungan", h.Create)
}