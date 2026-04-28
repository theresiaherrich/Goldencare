package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type LansiaRepository interface {
	Create(ctx context.Context, lansia *models.Lansia) error
	GetByID(ctx context.Context, id string) (*models.Lansia, error)
	GetByPantiID(ctx context.Context, pantiID string, filters map[string]interface{}) ([]models.Lansia, error)
	GetCountByPantiID(ctx context.Context, pantiID string) (int, error)
	Update(ctx context.Context, lansia *models.Lansia) error
}

type lansiaRepository struct {
	db *sqlx.DB
}

func newLansiaRepository(db *sqlx.DB) LansiaRepository {
	return &lansiaRepository{db: db}
}

func (r *lansiaRepository) Create(ctx context.Context, lansia *models.Lansia) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO lansia
		(id, panti_id, kamar_id, nama, nik, tanggal_lahir, jenis_kelamin,
		 alamat_asal, golongan_darah, riwayat_penyakit, alergi, foto_url, status, tanggal_masuk, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`, lansia.ID, lansia.PantiID, lansia.KamarID, lansia.Nama, lansia.NIK, lansia.TanggalLahir,
		lansia.JenisKelamin, lansia.AlamatAsal, lansia.GolonganDarah, lansia.RiwayatPenyakit,
		lansia.Alergi, lansia.FotoURL, lansia.Status, lansia.TanggalMasuk, lansia.CreatedAt, lansia.UpdatedAt)
	return err
}

func (r *lansiaRepository) GetByID(ctx context.Context, id string) (*models.Lansia, error) {
	lansia := &models.Lansia{}
	err := r.db.GetContext(ctx, lansia, `
		SELECT id, panti_id, kamar_id, nama, nik, tanggal_lahir, jenis_kelamin,
		       alamat_asal, golongan_darah, riwayat_penyakit, alergi, foto_url, status, tanggal_masuk, created_at, updated_at
		FROM lansia WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return lansia, nil
}

func (r *lansiaRepository) GetByPantiID(ctx context.Context, pantiID string, filters map[string]interface{}) ([]models.Lansia, error) {
	query := `
		SELECT id, panti_id, kamar_id, nama, nik, tanggal_lahir, jenis_kelamin,
		       alamat_asal, golongan_darah, riwayat_penyakit, alergi, foto_url, status, tanggal_masuk, created_at, updated_at
		FROM lansia WHERE panti_id = $1 AND status = 'aktif'
	`
	args := []interface{}{pantiID}
	argIdx := 2

	if kamar, ok := filters["kamar"].(string); ok && kamar != "" {
		query += fmt.Sprintf(" AND kamar_id = $%d", argIdx)
		args = append(args, kamar)
		argIdx++
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
	}
	query += " ORDER BY nama ASC"

	var lansias []models.Lansia
	if err := r.db.SelectContext(ctx, &lansias, query, args...); err != nil {
		return nil, err
	}
	return lansias, nil
}

func (r *lansiaRepository) GetCountByPantiID(ctx context.Context, pantiID string) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM lansia WHERE panti_id = $1 AND status = 'aktif'`, pantiID)
	return count, err
}

func (r *lansiaRepository) Update(ctx context.Context, lansia *models.Lansia) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE lansia
		SET panti_id = $1, kamar_id = $2, nama = $3, nik = $4, tanggal_lahir = $5,
		    jenis_kelamin = $6, alamat_asal = $7, golongan_darah = $8, riwayat_penyakit = $9,
		    alergi = $10, status = $11, tanggal_masuk = $12, foto_url = $13, updated_at = $14
		WHERE id = $15
	`, lansia.PantiID, lansia.KamarID, lansia.Nama, lansia.NIK, lansia.TanggalLahir,
		lansia.JenisKelamin, lansia.AlamatAsal, lansia.GolonganDarah, lansia.RiwayatPenyakit,
		lansia.Alergi, lansia.Status, lansia.TanggalMasuk, lansia.FotoURL, lansia.UpdatedAt, lansia.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("lansia with id %s not found", lansia.ID)
	}
	return nil
}