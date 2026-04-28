package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type CatatanShiftRepository interface {
	Create(ctx context.Context, catatan *models.CatatanShift) error
	GetByID(ctx context.Context, id string) (*models.CatatanShift, error)
	GetByLansia(ctx context.Context, lansiaID string, filters map[string]interface{}) ([]models.CatatanShift, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	GetDraftByPengurus(ctx context.Context, pengurusID string) ([]models.CatatanShift, error)
}

type catatanShiftRepository struct {
	db *sqlx.DB
}

func newCatatanShiftRepository(db *sqlx.DB) CatatanShiftRepository {
	return &catatanShiftRepository{db: db}
}

func (r *catatanShiftRepository) Create(ctx context.Context, catatan *models.CatatanShift) error {
	query := `
		INSERT INTO catatan_shift
		(id, lansia_id, pengurus_id, isi_catatan, shift,
		 suasana_hati, nafsu_makan, aktivitas, status_jurnal, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		catatan.ID, catatan.LansiaID, catatan.PengurusID, catatan.IsiCatatan, catatan.Shift,
		catatan.SuasanaHati, catatan.NafsuMakan, catatan.Aktivitas, catatan.StatusJurnal, catatan.CreatedAt,
	)
	return err
}

func (r *catatanShiftRepository) GetByID(ctx context.Context, id string) (*models.CatatanShift, error) {
	var catatan models.CatatanShift
	query := `SELECT id, lansia_id, pengurus_id, isi_catatan, shift,
	                 suasana_hati, nafsu_makan, aktivitas, status_jurnal, created_at
	          FROM catatan_shift WHERE id = $1`
	if err := r.db.GetContext(ctx, &catatan, query, id); err != nil {
		return nil, err
	}
	return &catatan, nil
}

func (r *catatanShiftRepository) GetByLansia(ctx context.Context, lansiaID string, filters map[string]interface{}) ([]models.CatatanShift, error) {
	query := `
		SELECT cs.*, u.name AS nama_pengurus
		FROM catatan_shift cs
		JOIN users u ON u.id = cs.pengurus_id
		WHERE cs.lansia_id = $1
	`
	args := []interface{}{lansiaID}
	argIdx := 2

	if status, ok := filters["status"].(string); ok && status != "" {
		query += fmt.Sprintf(" AND cs.status_jurnal = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if shift, ok := filters["shift"].(string); ok && shift != "" {
		query += fmt.Sprintf(" AND cs.shift = $%d", argIdx)
		args = append(args, shift)
		argIdx++
	}

	query += " ORDER BY cs.created_at DESC LIMIT 30"

	var catatans []models.CatatanShift
	if err := r.db.SelectContext(ctx, &catatans, query, args...); err != nil {
		return nil, err
	}
	return catatans, nil
}

func (r *catatanShiftRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE catatan_shift SET status_jurnal = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *catatanShiftRepository) GetDraftByPengurus(ctx context.Context, pengurusID string) ([]models.CatatanShift, error) {
	var drafts []models.CatatanShift
	query := `
		SELECT cs.*, l.nama AS nama_lansia
		FROM catatan_shift cs
		JOIN lansia l ON l.id = cs.lansia_id
		WHERE cs.pengurus_id = $1 AND cs.status_jurnal = 'draf'
		ORDER BY cs.created_at DESC
	`
	if err := r.db.SelectContext(ctx, &drafts, query, pengurusID); err != nil {
		return nil, err
	}
	return drafts, nil
}
