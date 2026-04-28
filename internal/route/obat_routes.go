package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/handlers"
)

func registerObatRoutes(protected fiber.Router, pengurus fiber.Router, h *handlers.ObatHandler) {
	protected.Get("/obat/options", h.GetOptions)

	pengurus.Post("/obat", h.Create)
	pengurus.Get("/obat/lansia/:lansia_id", h.GetByLansia)
	pengurus.Get("/obat/riwayat/:lansia_id", h.GetRiwayatPemberian)
	pengurus.Get("/obat/jadwal-hari-ini", h.GetJadwalHariIni)
	pengurus.Post("/obat/checklist/:jadwal_obat_id", h.Checklist)
	pengurus.Delete("/obat/:obat_id", h.DeleteObat)
}