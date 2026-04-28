package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

type PemeriksaanService interface {
	Create(ctx context.Context, req *models.CreatePemeriksaanRequest, pengurusID string) (*models.PemeriksaanKesehatan, *HealthRecommendation, error)
	GetByID(ctx context.Context, id string) (*models.PemeriksaanKesehatan, error)
	GetByLansia(ctx context.Context, lansiaID string) ([]models.PemeriksaanKesehatan, error)
	GetLatestWithAnalysis(ctx context.Context, lansiaID string) (*models.PemeriksaanKesehatan, *HealthRecommendation, error)
	AnalyzeOnly(ctx context.Context, req *models.CreatePemeriksaanRequest) (*HealthRecommendation, error)
}

type pemeriksaanService struct {
	repo repository.PemeriksaanRepository
}

func NewPemeriksaanService(repo repository.PemeriksaanRepository) PemeriksaanService {
	return &pemeriksaanService{repo: repo}
}

func (s *pemeriksaanService) Create(ctx context.Context, req *models.CreatePemeriksaanRequest, pengurusID string) (*models.PemeriksaanKesehatan, *HealthRecommendation, error) {
	if req.LansiaID == uuid.Nil {
		return nil, nil, errors.New("lansia_id wajib diisi")
	}

	recommendation := AnalyzeHealthCondition(*req)

	pemeriksaan := &models.PemeriksaanKesehatan{
		ID:            uuid.New(),
		LansiaID:      req.LansiaID,
		PengurusID:    uuid.MustParse(pengurusID),
		TekananDarah:  req.TekananDarah,
		DetakJantung:  req.DetakJantung,
		GulaDarah:     req.GulaDarah,
		SuhuTubuh:     req.SuhuTubuh,
		BeratBadan:    req.BeratBadan,
		Keluhan:       req.Keluhan,
		Status:        string(recommendation.Status),
		Rekomendasi:   recommendation.Description,
		StatusDarurat: string(recommendation.Color),
		CreatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, pemeriksaan); err != nil {
		return nil, nil, err
	}
	return pemeriksaan, &recommendation, nil
}

func (s *pemeriksaanService) GetByID(ctx context.Context, id string) (*models.PemeriksaanKesehatan, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *pemeriksaanService) GetByLansia(ctx context.Context, lansiaID string) ([]models.PemeriksaanKesehatan, error) {
	return s.repo.GetByLansia(ctx, lansiaID)
}

func (s *pemeriksaanService) GetLatestWithAnalysis(ctx context.Context, lansiaID string) (*models.PemeriksaanKesehatan, *HealthRecommendation, error) {
	pemeriksaan, err := s.repo.GetLatestByLansia(ctx, lansiaID)
	if err != nil {
		return nil, nil, err
	}
	req := models.CreatePemeriksaanRequest{
		LansiaID:     pemeriksaan.LansiaID,
		TekananDarah: pemeriksaan.TekananDarah,
		DetakJantung: pemeriksaan.DetakJantung,
		GulaDarah:    pemeriksaan.GulaDarah,
		SuhuTubuh:    pemeriksaan.SuhuTubuh,
		BeratBadan:   pemeriksaan.BeratBadan,
		Keluhan:      pemeriksaan.Keluhan,
	}
	recommendation := AnalyzeHealthCondition(req)
	return pemeriksaan, &recommendation, nil
}

func (s *pemeriksaanService) AnalyzeOnly(ctx context.Context, req *models.CreatePemeriksaanRequest) (*HealthRecommendation, error) {
	recommendation := AnalyzeHealthCondition(*req)
	return &recommendation, nil
}