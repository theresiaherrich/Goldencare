package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/handlers"
)

func registerCatatanShiftRoutes(
	protected fiber.Router,
	pengurus fiber.Router,
	h *handlers.CatatanShiftHandler,
) {
	protected.Get("/catatan-shift/options", h.GetOptions)

	pengurus.Post("/catatan-shift", h.Create)
	pengurus.Put("/catatan-shift/:id/kirim", h.KirimOperan)
	pengurus.Get("/catatan-shift/lansia/:lansia_id", h.GetByLansia)
	pengurus.Get("/catatan-shift/draf-saya", h.GetDraftSaya)
}