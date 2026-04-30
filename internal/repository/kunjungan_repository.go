package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type KunjunganRepository interface {
	Create(ctx context.Context, kunjungan *models.KunjunganKeluarga) error
	GetByLansia(ctx context.Context, lansiaID string) ([]models.KunjunganKeluarga, error)
	GetTerbaru(ctx context.Context, lansiaID string, limit int) ([]models.KunjunganKeluarga, error)
}

type kunjunganRepository struct {
	db *sqlx.DB
}

func newKunjunganRepository(db *sqlx.DB) KunjunganRepository {
	return &kunjunganRepository{db: db}
}

func (r *kunjunganRepository) Create(ctx context.Context, kunjungan *models.KunjunganKeluarga) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO kunjungan_keluarga
		(id, lansia_id, pengurus_id, nama_keluarga, hubungan_keluarga, tanggal_kunjungan, durasi_menit, foto_url, catatan, respon_lansia, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, kunjungan.ID, kunjungan.LansiaID, kunjungan.PengurusID, kunjungan.NamaKeluarga, kunjungan.HubunganKeluarga,
		kunjungan.TanggalKunjungan, kunjungan.DurasiMenit, kunjungan.FotoURL, kunjungan.Catatan,
		kunjungan.ResponLansia, kunjungan.CreatedAt)
	return err
}

func (r *kunjunganRepository) GetByLansia(ctx context.Context, lansiaID string) ([]models.KunjunganKeluarga, error) {
	var results []models.KunjunganKeluarga
	err := r.db.SelectContext(ctx, &results, `
    SELECT 
        kk.id,
        kk.lansia_id,
        kk.pengurus_id,
        kk.nama_keluarga,
        kk.hubungan_keluarga,
        kk.tanggal_kunjungan,
        kk.durasi_menit,
        COALESCE(kk.foto_url, '') AS foto_url,
        COALESCE(kk.catatan, '') AS catatan,
        COALESCE(kk.respon_lansia, '') AS respon_lansia,
        kk.created_at,
        u.name AS nama_pengurus
    FROM kunjungan_keluarga kk
    JOIN users u ON u.id = kk.pengurus_id
    WHERE kk.lansia_id = $1
    ORDER BY kk.tanggal_kunjungan DESC
`, lansiaID)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *kunjunganRepository) GetTerbaru(ctx context.Context, lansiaID string, limit int) ([]models.KunjunganKeluarga, error) {
	var results []models.KunjunganKeluarga
	err := r.db.SelectContext(ctx, &results, `
		SELECT 
			kk.id,
			kk.lansia_id,
			kk.pengurus_id,
			kk.nama_keluarga,
			kk.hubungan_keluarga,
			kk.tanggal_kunjungan,
			kk.durasi_menit,
			COALESCE(kk.foto_url, '') AS foto_url,
			COALESCE(kk.catatan, '') AS catatan,
			COALESCE(kk.respon_lansia, '') AS respon_lansia,
			kk.created_at,
			u.name AS nama_pengurus
		FROM kunjungan_keluarga kk
		JOIN users u ON u.id = kk.pengurus_id
		WHERE kk.lansia_id = $1
		ORDER BY kk.tanggal_kunjungan DESC
		LIMIT $2
	`, lansiaID, limit)
	if err != nil {
		return nil, err
	}
	return results, nil
}
