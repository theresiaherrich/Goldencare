package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/middleware"
	"github.com/theresiaherrich/Goldencare/internal/models"
	service "github.com/theresiaherrich/Goldencare/internal/services"
	"github.com/theresiaherrich/Goldencare/pkg/utils"
)

type ObatHandler struct {
	service service.ObatService
}

func NewObatHandler(svc service.ObatService) *ObatHandler {
	return &ObatHandler{service: svc}
}

func (h *ObatHandler) GetOptions(c *fiber.Ctx) error {
	return utils.OK(c, "Opsi kelola obat", fiber.Map{
		"cara_pemberian": models.CaraPemberianOptions,
		"shift":          models.ShiftObatOptions,
	})
}

func (h *ObatHandler) Create(c *fiber.Ctx) error {
	var req models.CreateObatRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Body tidak valid")
	}
	result, err := h.service.Create(c.Context(), &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.Created(c, "Obat dan jadwal berhasil ditambahkan", result)
}

func (h *ObatHandler) GetJadwalHariIni(c *fiber.Ctx) error {
	pengurusID := middleware.GetUserID(c)
	result, err := h.service.GetJadwalHariIni(c.Context(), pengurusID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil jadwal obat hari ini")
	}
	return utils.OK(c, "Jadwal obat hari ini", result)
}

func (h *ObatHandler) Checklist(c *fiber.Ctx) error {
	jadwalObatID := c.Params("jadwal_obat_id")
	pengurusID := middleware.GetUserID(c)

	var req struct {
		Catatan string `json:"catatan"`
	}
	c.BodyParser(&req)

	log, err := h.service.Checklist(c.Context(), jadwalObatID, pengurusID, req.Catatan)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.Created(c, "Pemberian obat berhasil dicatat", log)
}

func (h *ObatHandler) GetByLansia(c *fiber.Ctx) error {
	lansiaID := c.Params("lansia_id")
	result, err := h.service.GetByLansia(c.Context(), lansiaID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil data obat")
	}
	return utils.OK(c, "Daftar obat lansia", result)
}

func (h *ObatHandler) GetRiwayatPemberian(c *fiber.Ctx) error {
	lansiaID := c.Params("lansia_id")
	result, err := h.service.GetRiwayatPemberian(c.Context(), lansiaID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil riwayat pemberian obat")
	}
	return utils.OK(c, "Riwayat pemberian obat", result)
}

func (h *ObatHandler) DeleteObat(c *fiber.Ctx) error {
	obatID := c.Params("obat_id")
	if err := h.service.DeleteObat(c.Context(), obatID); err != nil {
		return utils.InternalError(c, "Gagal menonaktifkan obat")
	}
	return utils.OK(c, "Obat berhasil dinonaktifkan", nil)
}
