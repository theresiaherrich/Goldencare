package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type ObatRepository interface {
	CreateWithJadwal(ctx context.Context, obat *models.Obat, jadwals []models.JadwalObat) error
	GetByLansia(ctx context.Context, lansiaID string) ([]models.ObatDetail, error)
	GetRiwayatPemberian(ctx context.Context, lansiaID string) ([]models.RiwayatPemberianItem, error)
	GetJadwalHariIni(ctx context.Context, pengurusID string) ([]models.JadwalHariIniItem, error)
	CekSudahDiberikan(ctx context.Context, jadwalObatID string) (bool, error)
	CreateLog(ctx context.Context, log *models.LogPemberianObat) error
	SoftDelete(ctx context.Context, obatID string) error
}

type obatRepository struct {
	db *sqlx.DB
}

func newObatRepository(db *sqlx.DB) ObatRepository {
	return &obatRepository{db: db}
}

func (r *obatRepository) CreateWithJadwal(ctx context.Context, obat *models.Obat, jadwals []models.JadwalObat) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExecContext(ctx, `
		INSERT INTO obat (id, lansia_id, nama_obat, dosis, cara_pemberian, keterangan, is_aktif, created_at, updated_at)
		VALUES (:id, :lansia_id, :nama_obat, :dosis, :cara_pemberian, :keterangan, :is_aktif, :created_at, :updated_at)
	`, obat)
	if err != nil {
		return err
	}

	for _, j := range jadwals {
		_, err = tx.NamedExecContext(ctx, `
			INSERT INTO jadwal_obat (id, obat_id, jam, shift, created_at)
			VALUES (:id, :obat_id, :jam, :shift, :created_at)
		`, j)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *obatRepository) GetByLansia(ctx context.Context, lansiaID string) ([]models.ObatDetail, error) {
	var obats []models.Obat
	if err := r.db.SelectContext(ctx, &obats, `
		SELECT * FROM obat WHERE lansia_id = $1 AND is_aktif = true ORDER BY created_at DESC
	`, lansiaID); err != nil {
		return nil, err
	}

	result := make([]models.ObatDetail, 0, len(obats))
	for _, o := range obats {
		var jadwals []models.JadwalObat
		r.db.SelectContext(ctx, &jadwals, `SELECT * FROM jadwal_obat WHERE obat_id = $1 ORDER BY jam ASC`, o.ID)
		result = append(result, models.ObatDetail{Obat: o, Jadwals: jadwals})
	}
	return result, nil
}

func (r *obatRepository) GetRiwayatPemberian(ctx context.Context, lansiaID string) ([]models.RiwayatPemberianItem, error) {
	results := make([]models.RiwayatPemberianItem, 0)
	err := r.db.SelectContext(ctx, &results, `
		SELECT lp.*, o.nama_obat, o.dosis, o.cara_pemberian, jo.jam, jo.shift, u.name AS nama_pengurus
		FROM log_pemberian_obat lp
		JOIN jadwal_obat jo ON jo.id = lp.jadwal_obat_id
		JOIN obat o         ON o.id  = jo.obat_id
		JOIN users u        ON u.id  = lp.pengurus_id
		WHERE o.lansia_id = $1
		ORDER BY lp.diberikan_pada DESC
		LIMIT 50
	`, lansiaID)
	return results, err
}

func (r *obatRepository) GetJadwalHariIni(ctx context.Context, pengurusID string) ([]models.JadwalHariIniItem, error) {
	var items []models.JadwalHariIniItem
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			jo.id                       AS jadwal_obat_id,
			o.id                        AS obat_id,
			l.id                        AS lansia_id,
			l.nama                      AS nama_lansia,
			COALESCE(k.nama_kamar, '-') AS nomor_kamar,
			COALESCE(l.foto_url, '')    AS foto_lansia,
			o.nama_obat,
			o.dosis,
			o.cara_pemberian,
			jo.jam,
			jo.shift,
			lp.diberikan_pada,
			COALESCE(u.name, '')        AS nama_pengurus
		FROM jadwal_obat jo
		JOIN obat o        ON o.id = jo.obat_id
		JOIN lansia l      ON l.id = o.lansia_id
		LEFT JOIN kamar k  ON k.id = l.kamar_id
		LEFT JOIN log_pemberian_obat lp
			ON lp.jadwal_obat_id = jo.id AND DATE(lp.diberikan_pada) = CURRENT_DATE
		LEFT JOIN users u  ON u.id = lp.pengurus_id
		WHERE o.is_aktif = true
		  AND l.panti_id = (SELECT panti_id FROM users WHERE id = $1)
		ORDER BY jo.jam ASC
	`, pengurusID)
	return items, err
}

func (r *obatRepository) CekSudahDiberikan(ctx context.Context, jadwalObatID string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
		SELECT COUNT(*) FROM log_pemberian_obat
		WHERE jadwal_obat_id = $1 AND DATE(diberikan_pada) = CURRENT_DATE
	`, jadwalObatID)
	return count > 0, err
}

func (r *obatRepository) CreateLog(ctx context.Context, log *models.LogPemberianObat) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO log_pemberian_obat (id, jadwal_obat_id, pengurus_id, diberikan_pada, catatan)
		VALUES (:id, :jadwal_obat_id, :pengurus_id, :diberikan_pada, :catatan)
	`, log)
	return err
}

func (r *obatRepository) SoftDelete(ctx context.Context, obatID string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE obat SET is_aktif = false, updated_at = NOW() WHERE id = $1`, obatID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errNotFound("obat")
	}
	return nil
}