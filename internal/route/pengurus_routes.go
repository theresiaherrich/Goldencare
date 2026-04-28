package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/handlers"
)

func registerPengurusRoutes(
	protected fiber.Router,
	pengurus fiber.Router,
	galeriHandler *handlers.GaleriHandler,
	catatanShiftHandler *handlers.CatatanShiftHandler,
	pemeriksaanHandler *handlers.PemeriksaanHandler,
	kunjunganHandler *handlers.KunjunganHandler,
	obatHandler *handlers.ObatHandler,
	pengurusHandler *handlers.PengurusHandler,
) {
	registerGaleriRoutes(protected, pengurus, galeriHandler)
	registerCatatanShiftRoutes(protected, pengurus, catatanShiftHandler)
	registerPemeriksaanRoutes(protected, pengurus, pemeriksaanHandler)
	registerKunjunganRoutes(protected, pengurus, kunjunganHandler)
	registerObatRoutes(protected, pengurus, obatHandler)

	pengurus.Get("/dashboard", pengurusHandler.GetDashboard)
	pengurus.Get("", pengurusHandler.GetAll)
	pengurus.Get(":user_id", pengurusHandler.GetByID)
	pengurus.Post("/profil", pengurusHandler.SetProfil)
	pengurus.Get("/shift-saya", pengurusHandler.GetShiftSaya)
}