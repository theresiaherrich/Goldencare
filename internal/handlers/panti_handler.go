package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/theresiaherrich/Goldencare/internal/middleware"
	"github.com/theresiaherrich/Goldencare/internal/models"
	service "github.com/theresiaherrich/Goldencare/internal/services"
	"github.com/theresiaherrich/Goldencare/pkg/utils"
)

type PantiHandler struct {
	pantiService service.PantiService
}

func NewPantiHandler(pantiService service.PantiService) *PantiHandler {
	return &PantiHandler{pantiService: pantiService}
}

func (h *PantiHandler) GetPanti(c *fiber.Ctx) error {
	pantiID := middleware.GetPantiID(c)
	if pantiID == "" {
		return utils.BadRequest(c, "Anda belum terdaftar di panti")
	}
	panti, err := h.pantiService.GetByID(c.Context(), pantiID)
	if err != nil {
		return utils.NotFound(c, "Panti tidak ditemukan")
	}
	return utils.OK(c, "Data panti", panti)
}

func (h *PantiHandler) GetAll(c *fiber.Ctx) error {
	pantis, err := h.pantiService.GetAll(c.Context())
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil data panti")
	}
	return utils.OK(c, "Daftar panti", pantis)
}

func (h *PantiHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Nama    string `json:"nama"`
		Alamat  string `json:"alamat"`
		Telepon string `json:"telepon"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Request body tidak valid")
	}

	userID := middleware.GetUserID(c)
	panti := &models.Panti{
		ID:          uuid.New(),
		Nama:        req.Nama,
		Alamat:      &req.Alamat,
		Telepon:     &req.Telepon,
		PengelolaID: uuid.MustParse(userID),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := h.pantiService.Create(c.Context(), panti); err != nil {
		return utils.InternalError(c, "Gagal membuat panti")
	}
	return utils.Created(c, "Panti berhasil dibuat", fiber.Map{
		"id":   panti.ID,
		"nama": panti.Nama,
	})
}

func (h *PantiHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	panti, err := h.pantiService.GetByID(c.Context(), id)
	if err != nil {
		return utils.NotFound(c, "Panti tidak ditemukan")
	}

	var req struct {
		Nama    string `json:"nama"`
		Alamat  string `json:"alamat"`
		Telepon string `json:"telepon"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Request body tidak valid")
	}

	if req.Nama != "" {
		panti.Nama = req.Nama
	}
	if req.Alamat != "" {
		panti.Alamat = &req.Alamat
	}
	if req.Telepon != "" {
		panti.Telepon = &req.Telepon
	}
	panti.UpdatedAt = time.Now()

	if err := h.pantiService.Update(c.Context(), panti); err != nil {
		return utils.InternalError(c, "Gagal memperbarui panti")
	}
	return utils.OK(c, "Panti berhasil diperbarui", fiber.Map{
		"id":   panti.ID,
		"nama": panti.Nama,
	})
}
