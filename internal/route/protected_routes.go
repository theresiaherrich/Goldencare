package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/bootstrap"
	"github.com/theresiaherrich/Goldencare/internal/middleware"
)

func registerProtectedRoutes(api fiber.Router, deps *AppDependencies, container *bootstrap.Container) {
	protected := api.Group("", middleware.RequireAuth(container.Config))

	protected.Get("/auth/me", deps.AuthHandler.Me)

	pengelola := protected.Group("/pengelola", middleware.RequireRole("pengelola, superadmin"))
	registerPengelolaRoutes(
		pengelola,
		deps.PantiHandler,
		deps.AuthHandler,
		deps.LansiaHandler,
		deps.PengurusHandler,
	)

	pengurus := protected.Group("/pengurus", middleware.RequireRole("pengurus, superadmin"))
	registerPengurusRoutes(
		protected,
		pengurus,
		deps.GaleriHandler,
		deps.CatatanShiftHandler,
		deps.PemeriksaanHandler,
		deps.KunjunganHandler,
		deps.ObatHandler,
		deps.PengurusHandler,
	)

	keluarga := protected.Group("/keluarga", middleware.RequireRole("keluarga, superadmin"))
	registerKeluargaRoutes(
		keluarga,
		deps.KeluargaHandler,
	)
}
