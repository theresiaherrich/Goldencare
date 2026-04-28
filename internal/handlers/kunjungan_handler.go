package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/middleware"
	"github.com/theresiaherrich/Goldencare/internal/models"
	service "github.com/theresiaherrich/Goldencare/internal/services"
	"github.com/theresiaherrich/Goldencare/pkg/utils"
)

type KunjunganHandler struct {
	service service.KunjunganService
}

func NewKunjunganHandler(svc service.KunjunganService) *KunjunganHandler {
	return &KunjunganHandler{service: svc}
}

func (h *KunjunganHandler) GetByLansia(c *fiber.Ctx) error {
	lansiaID := c.Params("lansia_id")
	riwayat, err := h.service.GetByLansia(c.Context(), lansiaID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil riwayat kunjungan")
	}
	return utils.OK(c, "Riwayat kunjungan lansia", riwayat)
}

func (h *KunjunganHandler) GetTerbaru(c *fiber.Ctx) error {
	lansiaID := c.Params("lansia_id")
	riwayat, err := h.service.GetTerbaru(c.Context(), lansiaID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil kunjungan terbaru")
	}
	return utils.OK(c, "Kunjungan keluarga terbaru", riwayat)
}

func (h *KunjunganHandler) Create(c *fiber.Ctx) error {
	var req models.CreateKunjunganRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Body tidak valid")
	}

	pengurusID := middleware.GetUserID(c)
	if pengurusID == "" {
		return utils.Unauthorized(c, "Pengurus tidak terautentikasi")
	}

	kunjungan, err := h.service.Create(c.Context(), &req, pengurusID)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.Created(c, "Kunjungan keluarga berhasil dicatat", kunjungan)
}
