package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theresiaherrich/Goldencare/internal/middleware"
	"github.com/theresiaherrich/Goldencare/internal/models"
	service "github.com/theresiaherrich/Goldencare/internal/services"
	"github.com/theresiaherrich/Goldencare/pkg/utils"
)

type PemeriksaanHandler struct {
	service service.PemeriksaanService
}

func NewPemeriksaanHandler(service service.PemeriksaanService) *PemeriksaanHandler {
	return &PemeriksaanHandler{service: service}
}

func (h *PemeriksaanHandler) Create(c *fiber.Ctx) error {
	var req models.CreatePemeriksaanRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Body tidak valid")
	}

	pengurusID := middleware.GetUserID(c)

	pemeriksaan, recommendation, err := h.service.Create(c.Context(), &req, pengurusID)
	if err != nil {
		return utils.InternalError(c, "Gagal menyimpan hasil pemeriksaan: "+err.Error())
	}

	responseData := fiber.Map{
		"pemeriksaan":    pemeriksaan,
		"recommendation": recommendation,
	}

	return utils.Created(c, "Hasil pemeriksaan berhasil disimpan", responseData)
}

func (h *PemeriksaanHandler) GetByLansia(c *fiber.Ctx) error {
	lansiaID := c.Params("lansia_id")

	pemeriksaans, err := h.service.GetByLansia(c.Context(), lansiaID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil riwayat pemeriksaan: "+err.Error())
	}

	return utils.OK(c, "Riwayat pemeriksaan", pemeriksaans)
}

func (h *PemeriksaanHandler) GetLatestWithAnalysis(c *fiber.Ctx) error {
	lansiaID := c.Params("lansia_id")

	pemeriksaan, recommendation, err := h.service.GetLatestWithAnalysis(c.Context(), lansiaID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil pemeriksaan terbaru: "+err.Error())
	}

	responseData := fiber.Map{
		"pemeriksaan":    pemeriksaan,
		"recommendation": recommendation,
	}

	return utils.OK(c, "Pemeriksaan terbaru dengan analisis", responseData)
}

func (h *PemeriksaanHandler) AnalyzeOnly(c *fiber.Ctx) error {
	var req models.CreatePemeriksaanRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Body tidak valid")
	}

	recommendation, err := h.service.AnalyzeOnly(c.Context(), &req)
	if err != nil {
		return utils.InternalError(c, "Gagal menganalisis tanda vital: "+err.Error())
	}

	return utils.OK(c, "Hasil analisis tanda vital", recommendation)
}
