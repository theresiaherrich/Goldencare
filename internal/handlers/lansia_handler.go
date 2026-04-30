package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/theresiaherrich/Goldencare/internal/middleware"
	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/services"
	"github.com/theresiaherrich/Goldencare/pkg/utils"
)

type LansiaHandler struct {
	lansiaService service.LansiaService
}

func NewLansiaHandler(lansiaService service.LansiaService) *LansiaHandler {
	return &LansiaHandler{lansiaService: lansiaService}
}

func (h *LansiaHandler) GetAll(c *fiber.Ctx) error {
	pantiID := middleware.GetPantiID(c)

	if middleware.IsSuperadmin(c) {
		pantiID = c.Query("panti_id", pantiID)
	}

	if pantiID == "" {
		return utils.BadRequest(c, "Anda belum terdaftar di panti")
	}

	filters := make(map[string]interface{})
	if kamar := c.Query("kamar"); kamar != "" {
		filters["kamar"] = kamar
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	lansias, err := h.lansiaService.GetByPantiID(c.Context(), pantiID, filters)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil data lansia")
	}
	return utils.OK(c, "Daftar lansia", lansias)
}

func (h *LansiaHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	pantiID := middleware.GetPantiID(c)

	lansia, err := h.lansiaService.GetByID(c.Context(), id)
	if err != nil {
		return utils.NotFound(c, "Lansia tidak ditemukan")
	}

	if !middleware.IsSuperadmin(c) {
		if pantiID == "" {
			return utils.BadRequest(c, "Anda belum terdaftar di panti")
		}
		if lansia.PantiID.String() != pantiID {
			return utils.Forbidden(c, "Anda tidak memiliki akses ke lansia ini")
		}
	}

	return utils.OK(c, "Detail lansia", lansia)
}

func (h *LansiaHandler) GetDashboard(c *fiber.Ctx) error {
	pantiID := middleware.GetPantiID(c)

	if middleware.IsSuperadmin(c) {
		pantiID = c.Query("panti_id", pantiID)
	}

	if pantiID == "" {
		return utils.BadRequest(c, "Anda belum terdaftar di panti")
	}

	dashboard, err := h.lansiaService.GetDashboard(c.Context(), pantiID)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}
	return utils.OK(c, "Dashboard profil lansia", dashboard)
}

func (h *LansiaHandler) Create(c *fiber.Ctx) error {
	pantiID := middleware.GetPantiID(c)

	if middleware.IsSuperadmin(c) {
		pantiID = c.Query("panti_id", pantiID)
	}

	if pantiID == "" {
		return utils.BadRequest(c, "Anda belum terdaftar di panti")
	}

	var req struct {
		Nama            string `json:"nama"`
		NIK             string `json:"nik"`
		TanggalLahir    string `json:"tanggal_lahir"`
		JenisKelamin    string `json:"jenis_kelamin"`
		AlamatAsal      string `json:"alamat_asal"`
		GolonganDarah   string `json:"golongan_darah"`
		RiwayatPenyakit string `json:"riwayat_penyakit"`
		Alergi          string `json:"alergi"`
		TanggalMasuk    string `json:"tanggal_masuk"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Request body tidak valid")
	}

	var tanggalLahir *time.Time
	if req.TanggalLahir != "" {
		t, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err != nil {
			return utils.BadRequest(c, "Format tanggal lahir tidak valid")
		}
		tanggalLahir = &t
	}

	tanggalMasuk := time.Now()
	if req.TanggalMasuk != "" {
		t, err := time.Parse("2006-01-02", req.TanggalMasuk)
		if err != nil {
			return utils.BadRequest(c, "Format tanggal masuk tidak valid")
		}
		tanggalMasuk = t
	}

	lansia := &models.Lansia{
		ID:              uuid.New(),
		PantiID:         uuid.MustParse(pantiID),
		Nama:            req.Nama,
		NIK:             req.NIK,
		TanggalLahir:    tanggalLahir,
		JenisKelamin:    req.JenisKelamin,
		AlamatAsal:      req.AlamatAsal,
		GolonganDarah:   req.GolonganDarah,
		RiwayatPenyakit: req.RiwayatPenyakit,
		Alergi:          req.Alergi,
		Status:          "aktif",
		TanggalMasuk:    tanggalMasuk,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.lansiaService.Create(c.Context(), lansia); err != nil {
		return utils.InternalError(c, "Gagal membuat lansia baru")
	}
	return utils.Created(c, "Lansia berhasil ditambahkan", fiber.Map{
		"id":   lansia.ID,
		"nama": lansia.Nama,
	})
}

func (h *LansiaHandler) Update(c *fiber.Ctx) error {
	pantiID := middleware.GetPantiID(c)

	if middleware.IsSuperadmin(c) {
		pantiID = c.Query("panti_id", pantiID)
	}

	if pantiID == "" {
		return utils.BadRequest(c, "Anda belum terdaftar di panti")
	}

	id := c.Params("id")
	lansia, err := h.lansiaService.GetByID(c.Context(), id)
	if err != nil {
		return utils.NotFound(c, "Lansia tidak ditemukan")
	}
	if lansia.PantiID.String() != pantiID {
		return utils.Forbidden(c, "Anda tidak memiliki akses ke lansia ini")
	}

	var req struct {
		Nama            string `json:"nama"`
		NIK             string `json:"nik"`
		TanggalLahir    string `json:"tanggal_lahir"`
		JenisKelamin    string `json:"jenis_kelamin"`
		AlamatAsal      string `json:"alamat_asal"`
		GolonganDarah   string `json:"golongan_darah"`
		RiwayatPenyakit string `json:"riwayat_penyakit"`
		Alergi          string `json:"alergi"`
		Status          string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Request body tidak valid")
	}

	if req.Nama != "" {
		lansia.Nama = req.Nama
	}
	if req.NIK != "" {
		lansia.NIK = req.NIK
	}
	if req.TanggalLahir != "" {
		if t, err := time.Parse("2006-01-02", req.TanggalLahir); err == nil {
			lansia.TanggalLahir = &t
		}
	}
	if req.JenisKelamin != "" {
		lansia.JenisKelamin = req.JenisKelamin
	}
	if req.AlamatAsal != "" {
		lansia.AlamatAsal = req.AlamatAsal
	}
	if req.GolonganDarah != "" {
		lansia.GolonganDarah = req.GolonganDarah
	}
	if req.RiwayatPenyakit != "" {
		lansia.RiwayatPenyakit = req.RiwayatPenyakit
	}
	if req.Alergi != "" {
		lansia.Alergi = req.Alergi
	}
	if req.Status != "" {
		lansia.Status = req.Status
	}
	lansia.UpdatedAt = time.Now()

	if err := h.lansiaService.Update(c.Context(), lansia); err != nil {
		return utils.InternalError(c, "Gagal memperbarui lansia")
	}
	return utils.OK(c, "Lansia berhasil diperbarui", fiber.Map{
		"id":   lansia.ID,
		"nama": lansia.Nama,
	})
}
