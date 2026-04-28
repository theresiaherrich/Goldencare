package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/middleware"
	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/services"
	"github.com/theresiaherrich/Goldencare/pkg/utils"
)

type CatatanShiftHandler struct {
	service service.CatatanShiftService
}

func NewCatatanShiftHandler(service service.CatatanShiftService) *CatatanShiftHandler {
	return &CatatanShiftHandler{service: service}
}

func (h *CatatanShiftHandler) GetOptions(c *fiber.Ctx) error {
	return utils.OK(c, "Opsi jurnal jaga", fiber.Map{
		"suasana_hati": models.SuasanaHatiOptions,
		"nafsu_makan":  models.NafsuMakanOptions,
		"aktivitas":    models.AktivitasOptions,
		"shift":        models.ShiftOptions,
	})
}

func (h *CatatanShiftHandler) Create(c *fiber.Ctx) error {
	var req models.CreateCatatanShiftRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Body tidak valid")
	}

	pengurusID := middleware.GetUserID(c)

	catatan, err := h.service.Create(c.Context(), &req, pengurusID)
	if err != nil {
		return utils.InternalError(c, "Gagal menyimpan jurnal jaga: "+err.Error())
	}

	return utils.Created(c, "Jurnal jaga berhasil disimpan", catatan)
}

func (h *CatatanShiftHandler) KirimOperan(c *fiber.Ctx) error {
	id := c.Params("id")
	pengurusID := middleware.GetUserID(c)

	if err := h.service.UpdateStatus(c.Context(), id, "terkirim", pengurusID); err != nil {
		return utils.InternalError(c, "Gagal mengirim operan jaga: "+err.Error())
	}

	catatan, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil data catatan")
	}

	return utils.OK(c, "Operan jaga berhasil dikirim", catatan)
}

func (h *CatatanShiftHandler) GetByLansia(c *fiber.Ctx) error {
	lansiaID := c.Params("lansia_id")
	filters := map[string]interface{}{
		"status": c.Query("status"),
		"shift":  c.Query("shift"),
	}

	catatans, err := h.service.GetByLansia(c.Context(), lansiaID, filters)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil jurnal jaga: "+err.Error())
	}

	return utils.OK(c, "Jurnal jaga lansia", catatans)
}

func (h *CatatanShiftHandler) GetDraftSaya(c *fiber.Ctx) error {
	pengurusID := middleware.GetUserID(c)

	drafts, err := h.service.GetDraftByPengurus(c.Context(), pengurusID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil draf jurnal: "+err.Error())
	}

	return utils.OK(c, "Draf jurnal jaga saya", drafts)
}
