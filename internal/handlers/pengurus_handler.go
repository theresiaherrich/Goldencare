package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/middleware"
	"github.com/theresiaherrich/Goldencare/internal/models"
	service "github.com/theresiaherrich/Goldencare/internal/services"
	"github.com/theresiaherrich/Goldencare/pkg/utils"
)

type PengurusHandler struct {
	service service.PengurusService
}

func NewPengurusHandler(svc service.PengurusService) *PengurusHandler {
	return &PengurusHandler{service: svc}
}

func (h *PengurusHandler) GetDashboard(c *fiber.Ctx) error {
	pantiID := middleware.GetPantiID(c)
	result, err := h.service.GetDashboard(c.Context(), pantiID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil dashboard")
	}
	return utils.OK(c, "Dashboard manajemen staf", result)
}

func (h *PengurusHandler) GetAll(c *fiber.Ctx) error {
	pantiID := middleware.GetPantiID(c)
	filterShift := c.Query("shift")
	result, err := h.service.GetAll(c.Context(), pantiID, filterShift)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil data staf")
	}
	return utils.OK(c, "Daftar staf", result)
}

func (h *PengurusHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	pantiID := middleware.GetPantiID(c)
	result, err := h.service.GetByID(c.Context(), userID, pantiID)
	if err != nil {
		return utils.NotFound(c, "Staf tidak ditemukan")
	}
	return utils.OK(c, "Detail staf", result)
}

func (h *PengurusHandler) SetProfil(c *fiber.Ctx) error {
	var req models.CreatePengurusProfilRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Body tidak valid")
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)

	targetUserID := userID
	if role == "pengelola" {
		if tid := c.Query("user_id"); tid != "" {
			targetUserID = tid
		}
	}

	if err := h.service.SetProfil(c.Context(), &req, targetUserID); err != nil {
		return utils.InternalError(c, "Gagal menyimpan profil pengurus")
	}
	return utils.Created(c, "Profil staf berhasil disimpan", nil)
}

func (h *PengurusHandler) GetShiftSaya(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	shifts, err := h.service.GetShiftSaya(c.Context(), userID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil jadwal shift")
	}
	return utils.OK(c, "Jadwal shift saya", shifts)
}

func (h *PengurusHandler) GetKamar(c *fiber.Ctx) error {
	pantiID := middleware.GetPantiID(c)
	result, err := h.service.GetKamar(c.Context(), pantiID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil data kamar")
	}
	return utils.OK(c, "Daftar kamar", result)
}

func (h *PengurusHandler) CreateKamar(c *fiber.Ctx) error {
	pantiID := middleware.GetPantiID(c)

	var req struct {
		NamaKamar string `json:"nama_kamar"`
		Kapasitas int    `json:"kapasitas"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Body tidak valid")
	}
	if req.Kapasitas == 0 {
		req.Kapasitas = 1
	}

	kamar, err := h.service.CreateKamar(c.Context(), pantiID, req.NamaKamar, req.Kapasitas)
	if err != nil {
		return utils.InternalError(c, "Gagal membuat kamar")
	}
	return utils.Created(c, "Kamar berhasil ditambahkan", kamar)
}
