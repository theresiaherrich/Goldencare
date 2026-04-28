package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

type CatatanShiftService interface {
	Create(ctx context.Context, req *models.CreateCatatanShiftRequest, pengurusID string) (*models.CatatanShift, error)
	GetByID(ctx context.Context, id string) (*models.CatatanShift, error)
	GetByLansia(ctx context.Context, lansiaID string, filters map[string]interface{}) ([]models.CatatanShift, error)
	UpdateStatus(ctx context.Context, id string, status string, pengurusID string) error
	GetDraftByPengurus(ctx context.Context, pengurusID string) ([]models.CatatanShift, error)
}

type catatanShiftService struct {
	repo repository.CatatanShiftRepository
}

func NewCatatanShiftService(repo repository.CatatanShiftRepository) CatatanShiftService {
	return &catatanShiftService{repo: repo}
}

func (s *catatanShiftService) Create(ctx context.Context, req *models.CreateCatatanShiftRequest, pengurusID string) (*models.CatatanShift, error) {
	if req.LansiaID == uuid.Nil {
		return nil, errors.New("lansia_id wajib diisi")
	}
	if req.IsiCatatan == "" {
		return nil, errors.New("isi_catatan wajib diisi")
	}
	if req.Shift == "" {
		return nil, errors.New("shift wajib diisi")
	}

	status := req.StatusJurnal
	if status == "" {
		status = "draf"
	}

	catatan := &models.CatatanShift{
		ID:           uuid.New(),
		LansiaID:     req.LansiaID,
		PengurusID:   uuid.MustParse(pengurusID),
		IsiCatatan:   req.IsiCatatan,
		Shift:        req.Shift,
		SuasanaHati:  req.SuasanaHati,
		NafsuMakan:   req.NafsuMakan,
		Aktivitas:    req.Aktivitas,
		StatusJurnal: status,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, catatan); err != nil {
		return nil, err
	}
	return catatan, nil
}

func (s *catatanShiftService) GetByID(ctx context.Context, id string) (*models.CatatanShift, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *catatanShiftService) GetByLansia(ctx context.Context, lansiaID string, filters map[string]interface{}) ([]models.CatatanShift, error) {
	return s.repo.GetByLansia(ctx, lansiaID, filters)
}

func (s *catatanShiftService) UpdateStatus(ctx context.Context, id string, status string, pengurusID string) error {
	catatan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if catatan.PengurusID.String() != pengurusID {
		return errors.New("tidak berhak mengubah catatan ini")
	}
	if catatan.StatusJurnal == "terkirim" {
		return errors.New("catatan sudah terkirim sebelumnya")
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *catatanShiftService) GetDraftByPengurus(ctx context.Context, pengurusID string) ([]models.CatatanShift, error) {
	return s.repo.GetDraftByPengurus(ctx, pengurusID)
}