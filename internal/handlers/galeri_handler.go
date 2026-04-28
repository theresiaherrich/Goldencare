package handlers

import (
	"io"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/theresiaherrich/Goldencare/internal/middleware"
	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/services"
	"github.com/theresiaherrich/Goldencare/pkg/utils"
)

type GaleriHandler struct {
	galeriService service.GaleriService
}

func NewGaleriHandler(galeriService service.GaleriService) *GaleriHandler {
	return &GaleriHandler{galeriService: galeriService}
}

func (h *GaleriHandler) Create(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return utils.BadRequest(c, "Gagal parse form")
	}

	files := form.File["foto"]
	if len(files) == 0 {
		return utils.BadRequest(c, "Field 'foto' wajib diisi (upload file JPG/PNG)")
	}
	file := files[0]

	lansiaIDStr := c.FormValue("lansia_id")
	if lansiaIDStr == "" {
		return utils.BadRequest(c, "lansia_id wajib diisi")
	}

	lansiaID, err := uuid.Parse(lansiaIDStr)
	if err != nil {
		return utils.BadRequest(c, "Format lansia_id tidak valid (harus UUID)")
	}

	src, err := file.Open()
	if err != nil {
		return utils.InternalError(c, "Gagal baca file")
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		return utils.InternalError(c, "Gagal baca data file")
	}

	req := &models.CreateGaleriRequest{
		LansiaID:   lansiaID,
		LokasiLuka: c.FormValue("lokasi_luka"),
		Deskripsi:  c.FormValue("deskripsi"),
	}

	pengurusID := middleware.GetUserID(c)
	contentType := file.Header.Get("Content-Type")

	galeri, err := h.galeriService.Create(c.Context(), req, pengurusID, file.Filename, contentType, fileData)
	if err != nil {
		log.Printf("Error create galeri: %v", err)
		return utils.InternalError(c, "Gagal menyimpan data galeri fisik")
	}

	return utils.Created(c, "Dokumentasi fisik berhasil disimpan", galeri)
}

func (h *GaleriHandler) GetOptions(c *fiber.Ctx) error {
	return utils.OK(c, "Opsi galeri fisik", fiber.Map{
		"lokasi_tubuh": []string{
			"Kepala", "Wajah", "Leher",
			"Bahu Kanan", "Bahu Kiri",
			"Lengan Kanan", "Lengan Kiri",
			"Tangan Kanan", "Tangan Kiri",
			"Dada", "Punggung", "Perut",
			"Pinggul", "Paha Kanan", "Paha Kiri",
			"Lutut Kanan", "Lutut Kiri",
			"Kaki Kanan", "Kaki Kiri",
			"Tumit Kanan", "Tumit Kiri",
		},
		"jenis_kondisi": []string{
			"Memar", "Luka Lecet", "Luka Sayat", "Luka Bakar",
			"Kemerahan", "Bengkak", "Ruam Kulit", "Infeksi",
			"Dekubitus (Luka Tekan)", "Bekas Jatuh", "Lainnya",
		},
	})
}

func (h *GaleriHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	galeri, err := h.galeriService.GetByID(c.Context(), id)
	if err != nil {
		return utils.NotFound(c, "Data tidak ditemukan")
	}
	return utils.OK(c, "Detail kondisi fisik", galeri)
}

func (h *GaleriHandler) GetByLansia(c *fiber.Ctx) error {
	lansiaID := c.Params("lansia_id")
	galeris, err := h.galeriService.GetByLansia(c.Context(), lansiaID)
	if err != nil {
		return utils.InternalError(c, "Gagal mengambil galeri fisik lansia")
	}
	return utils.OK(c, "Galeri fisik lansia", galeris)
}
