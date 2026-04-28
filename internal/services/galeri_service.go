package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/theresiaherrich/Goldencare/internal/config"
	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

type GaleriService interface {
	Create(ctx context.Context, req *models.CreateGaleriRequest, pengurusID, fileName, contentType string, fileData []byte) (*models.GaleriFisik, error)
	GetByID(ctx context.Context, id string) (*models.GaleriFisik, error)
	GetByLansia(ctx context.Context, lansiaID string) ([]models.GaleriFisik, error)
}

type galeriService struct {
	repo     repository.GaleriRepository
	supabase *SupabaseService
}

func NewGaleriService(repo repository.GaleriRepository, cfg *config.Config) GaleriService {
	svc := &galeriService{repo: repo}
	if cfg != nil {
		svc.supabase = NewSupabaseService(cfg)
	}
	return svc
}

func (s *galeriService) Create(ctx context.Context, req *models.CreateGaleriRequest, pengurusID, fileName, contentType string, fileData []byte) (*models.GaleriFisik, error) {
	if s.supabase == nil {
		return nil, errors.New("supabase configuration is not available")
	}
	if req.LansiaID == uuid.Nil {
		return nil, errors.New("lansia_id wajib diisi")
	}

	storagePath := fmt.Sprintf("galeri/%s/%d-%s", req.LansiaID.String(), time.Now().UnixNano(), fileName)
	fotoURL, err := s.supabase.UploadFile("galeri-fisik", storagePath, fileData, contentType)
	if err != nil {
		return nil, err
	}

	tingkat, prediksi := analyzeGaleriCondition(req.Deskripsi, req.LokasiLuka)
	galeri := &models.GaleriFisik{
		ID:               uuid.New(),
		LansiaID:         req.LansiaID,
		PengurusID:       uuid.MustParse(pengurusID),
		FotoURL:          fotoURL,
		LokasiLuka:       req.LokasiLuka,
		Deskripsi:        req.Deskripsi,
		AnalisisAI:       "Analisis berdasarkan deskripsi: " + prediksi,
		TingkatDarurat:   tingkat,
		PrediksiPenyakit: prediksi,
		CreatedAt:        time.Now(),
	}

	if err := s.repo.Create(ctx, galeri); err != nil {
		return nil, err
	}
	return galeri, nil
}

func (s *galeriService) GetByID(ctx context.Context, id string) (*models.GaleriFisik, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *galeriService) GetByLansia(ctx context.Context, lansiaID string) ([]models.GaleriFisik, error) {
	return s.repo.GetByLansia(ctx, lansiaID)
}

func analyzeGaleriCondition(deskripsi, lokasi string) (string, string) {
	kataKritikal := []string{"berdarah banyak", "tidak sadar", "bengkak parah", "patah", "sesak"}
	kataSedang := []string{"memar", "luka dalam", "bernanah", "infeksi", "bengkak"}

	for _, k := range kataKritikal {
		if strings.Contains(strings.ToLower(deskripsi), strings.ToLower(k)) {
			return "kritis", "⚠️ Kondisi kritis! Segera hubungi dokter atau bawa ke IGD."
		}
	}
	for _, k := range kataSedang {
		if strings.Contains(strings.ToLower(deskripsi), strings.ToLower(k)) {
			return "sedang", "Kondisi memerlukan perhatian medis. Konsultasikan dengan dokter panti dalam 24 jam."
		}
	}
	return "ringan", "Kondisi luka ringan. Bersihkan dengan antiseptik dan pantau perkembangan."
}