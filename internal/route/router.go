package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/bootstrap"
	"github.com/theresiaherrich/Goldencare/internal/handlers"
	"github.com/theresiaherrich/Goldencare/internal/middleware"
)

type AppDependencies struct {
	AuthHandler         *handlers.AuthHandler
	PantiHandler        *handlers.PantiHandler
	LansiaHandler       *handlers.LansiaHandler
	PengurusHandler     *handlers.PengurusHandler
	KunjunganHandler    *handlers.KunjunganHandler
	CatatanShiftHandler *handlers.CatatanShiftHandler
	KeluargaHandler     *handlers.KeluargaHandler
	ObatHandler         *handlers.ObatHandler
	PemeriksaanHandler  *handlers.PemeriksaanHandler
	GaleriHandler       *handlers.GaleriHandler
}

func RegisterRoutes(app *fiber.App, container *bootstrap.Container) {
	app.Use(middleware.Logger())

	deps := buildDependencies(container)

	api := app.Group("/api")

	registerHealthRoutes(api, container)
	registerPublicRoutes(api, deps)
	registerSuperadminRoutes(api, deps, container)
	registerProtectedRoutes(api, deps, container)
}

func buildDependencies(container *bootstrap.Container) *AppDependencies {
	return &AppDependencies{
		AuthHandler:        handlers.NewAuthHandler(container.Service.Auth()),
		PantiHandler:       handlers.NewPantiHandler(container.Service.Panti()),
		LansiaHandler:      handlers.NewLansiaHandler(container.Service.Lansia()),
		PengurusHandler:    handlers.NewPengurusHandler(container.Service.Pengurus()),
		KunjunganHandler:   handlers.NewKunjunganHandler(container.Service.Kunjungan()),
		KeluargaHandler:    handlers.NewKeluargaHandler(container.Service.Keluarga()),
		ObatHandler:        handlers.NewObatHandler(container.Service.Obat()),
		PemeriksaanHandler: handlers.NewPemeriksaanHandler(container.Service.Pemeriksaan()),
		GaleriHandler:      handlers.NewGaleriHandler(container.Service.Galeri()),
	}
}
