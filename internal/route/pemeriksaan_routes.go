package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/handlers"
)

func registerPemeriksaanRoutes(
	protected fiber.Router,
	pengurus fiber.Router,
	h *handlers.PemeriksaanHandler,
) {
	protected.Get("/pemeriksaan/lansia/:lansia_id", h.GetByLansia)
	protected.Get("/pemeriksaan/lansia/:lansia_id/terbaru", h.GetLatestWithAnalysis)

	pengurus.Post("/pemeriksaan", h.Create)
	pengurus.Post("/pemeriksaan/analyze", h.AnalyzeOnly)
}