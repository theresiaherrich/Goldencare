package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type PengurusRepository interface {
	GetDashboard(ctx context.Context, pantiID string) (*models.PengurusDashboard, error)
	GetAll(ctx context.Context, pantiID string, filterShift string) ([]models.PengurusDetail, error)
	GetByID(ctx context.Context, userID string, pantiID string) (*models.PengurusDetailResponse, error)
	SetProfil(ctx context.Context, profil *models.PengurusProfil, shifts []models.Shift) error
	GetShiftByUser(ctx context.Context, userID string) ([]models.Shift, error)
	GetKamar(ctx context.Context, pantiID string) ([]models.KamarWithStaf, error)
	CreateKamar(ctx context.Context, kamar *models.Kamar) error
}

type pengurusRepository struct {
	db *sqlx.DB
}

func newPengurusRepository(db *sqlx.DB) PengurusRepository {
	return &pengurusRepository{db: db}
}

func (r *pengurusRepository) GetDashboard(ctx context.Context, pantiID string) (*models.PengurusDashboard, error) {
	hariIni := hariIndonesia(time.Now().Weekday())

	var totalStaf int
	r.db.GetContext(ctx, &totalStaf, `SELECT COUNT(*) FROM users WHERE panti_id = $1 AND role = 'pengurus'`, pantiID)

	var stafTanpaShift int
	r.db.GetContext(ctx, &stafTanpaShift, `
		SELECT COUNT(*) FROM users u
		WHERE u.panti_id = $1 AND u.role = 'pengurus'
		  AND u.id NOT IN (SELECT DISTINCT s.user_id FROM shift s WHERE s.hari LIKE '%' || $2 || '%')
	`, pantiID, hariIni)

	var shiftsHariIni []models.ShiftHariIni
	r.db.SelectContext(ctx, &shiftsHariIni, `
		SELECT s.nama_shift, s.jam_mulai::TEXT, s.jam_selesai::TEXT, COUNT(s.user_id) AS jumlah_staf
		FROM shift s
		JOIN users u ON u.id = s.user_id
		WHERE u.panti_id = $1 AND s.hari LIKE '%' || $2 || '%'
		GROUP BY s.nama_shift, s.jam_mulai, s.jam_selesai
		ORDER BY s.jam_mulai
	`, pantiID, hariIni)

	return &models.PengurusDashboard{
		TotalStaf:      totalStaf,
		StafTanpaShift: stafTanpaShift,
		ShiftHariIni:   shiftsHariIni,
		Hari:           hariIni,
	}, nil
}

func (r *pengurusRepository) GetAll(ctx context.Context, pantiID string, filterShift string) ([]models.PengurusDetail, error) {
	hariIni := hariIndonesia(time.Now().Weekday())
	query := `
		SELECT
			u.id AS user_id, u.name AS nama, u.email,
			COALESCE(pp.jabatan, '') AS jabatan,
			COALESCE(k.nama_kamar, 'Belum ditentukan') AS nama_kamar,
			COALESCE(s.nama_shift, 'Tidak ada shift') AS nama_shift,
			COALESCE(s.jam_mulai::TEXT, '') AS jam_mulai,
			COALESCE(s.jam_selesai::TEXT, '') AS jam_selesai
		FROM users u
		LEFT JOIN pengurus_profil pp ON pp.user_id = u.id
		LEFT JOIN kamar k ON k.id = pp.kamar_id
		LEFT JOIN shift s ON s.user_id = u.id AND s.hari LIKE '%' || $2 || '%'
		WHERE u.panti_id = $1 AND u.role = 'pengurus'
	`
	args := []interface{}{pantiID, hariIni}
	if filterShift != "" {
		args = append(args, "%"+filterShift+"%")
		query += " AND LOWER(s.nama_shift) LIKE LOWER($3)"
	}
	query += " ORDER BY u.name ASC"

	var results []models.PengurusDetail
	if err := r.db.SelectContext(ctx, &results, query, args...); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *pengurusRepository) GetByID(ctx context.Context, userID string, pantiID string) (*models.PengurusDetailResponse, error) {
	var user models.User
	if err := r.db.GetContext(ctx, &user,
		`SELECT * FROM users WHERE id = $1 AND panti_id = $2 AND role = 'pengurus'`, userID, pantiID); err != nil {
		return nil, err
	}

	var profil models.PengurusProfil
	r.db.GetContext(ctx, &profil, `SELECT * FROM pengurus_profil WHERE user_id = $1`, userID)

	var shifts []models.Shift
	r.db.SelectContext(ctx, &shifts, `SELECT * FROM shift WHERE user_id = $1 ORDER BY jam_mulai`, userID)

	var namaKamar string
	if profil.KamarID != nil {
		r.db.GetContext(ctx, &namaKamar, `SELECT nama_kamar FROM kamar WHERE id = $1`, profil.KamarID)
	}

	return &models.PengurusDetailResponse{
		User:      user,
		Profil:    profil,
		Shifts:    shifts,
		NamaKamar: namaKamar,
	}, nil
}

func (r *pengurusRepository) SetProfil(ctx context.Context, profil *models.PengurusProfil, shifts []models.Shift) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.ExecContext(ctx, `DELETE FROM pengurus_profil WHERE user_id = $1`, profil.UserID)
	if _, err = tx.NamedExecContext(ctx, `
		INSERT INTO pengurus_profil (id, user_id, kamar_id, jabatan, created_at, updated_at)
		VALUES (:id, :user_id, :kamar_id, :jabatan, :created_at, :updated_at)
	`, profil); err != nil {
		return err
	}

	tx.ExecContext(ctx, `DELETE FROM shift WHERE user_id = $1`, profil.UserID)
	for _, s := range shifts {
		if _, err = tx.NamedExecContext(ctx, `
			INSERT INTO shift (id, user_id, nama_shift, jam_mulai, jam_selesai, hari, created_at)
			VALUES (:id, :user_id, :nama_shift, :jam_mulai, :jam_selesai, :hari, :created_at)
		`, s); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *pengurusRepository) GetShiftByUser(ctx context.Context, userID string) ([]models.Shift, error) {
	var shifts []models.Shift
	err := r.db.SelectContext(ctx, &shifts, `SELECT * FROM shift WHERE user_id = $1 ORDER BY jam_mulai`, userID)
	return shifts, err
}

func (r *pengurusRepository) GetKamar(ctx context.Context, pantiID string) ([]models.KamarWithStaf, error) {
	var results []models.KamarWithStaf
	err := r.db.SelectContext(ctx, &results, `
		SELECT k.*,
			(SELECT COUNT(*) FROM pengurus_profil pp WHERE pp.kamar_id = k.id) AS jumlah_staf,
			(SELECT COUNT(*) FROM lansia l WHERE l.kamar_id = k.id AND l.status = 'aktif') AS jumlah_lansia
		FROM kamar k
		WHERE k.panti_id = $1
		ORDER BY k.nama_kamar
	`, pantiID)
	return results, err
}

func (r *pengurusRepository) CreateKamar(ctx context.Context, kamar *models.Kamar) error {
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO kamar (id, panti_id, nama_kamar, kapasitas, created_at)
		VALUES (:id, :panti_id, :nama_kamar, :kapasitas, :created_at)
	`, kamar)
	return err
}