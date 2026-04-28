package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/handlers"
)

func registerGaleriRoutes(
	protected fiber.Router,
	pengurus fiber.Router,
	h *handlers.GaleriHandler,
) {
	protected.Get("/galeri/options", h.GetOptions)
	protected.Get("/galeri/:id", h.GetByID)
	protected.Get("/galeri/lansia/:lansia_id", h.GetByLansia)
	pengurus.Post("/galeri", h.Create)
}
