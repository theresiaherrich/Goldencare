package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/services"
	"github.com/theresiaherrich/Goldencare/pkg/utils"
)

type KeluargaHandler struct {
	service service.KeluargaService
}

func NewKeluargaHandler(svc service.KeluargaService) *KeluargaHandler {
	return &KeluargaHandler{service: svc}
}

func (h *KeluargaHandler) GetDashboard(c *fiber.Ctx) error {
	lansiaID := c.Params("lansia_id")
	result, err := h.service.GetDashboard(c.Context(), lansiaID)
	if err != nil {
		return utils.NotFound(c, "Lansia tidak ditemukan")
	}
	return utils.OK(c, "Dashboard perawatan harian", result)
}

func (h *KeluargaHandler) GetKeluarga(c *fiber.Ctx) error {
	lansiaID := c.Params("lansia_id")
	result, err := h.service.GetKeluarga(c.Context(), lansiaID)
	if err != nil {
		return utils.NotFound(c, "Lansia tidak ditemukan")
	}
	return utils.OK(c, "Keluarga perawatan lengkap", result)
}
