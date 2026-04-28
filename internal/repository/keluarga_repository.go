package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/models"
)

type KeluargaRepository interface {
	GetDashboard(ctx context.Context, lansiaID string) (*models.KeluargaDashboard, error)
	GetKeluarga(ctx context.Context, lansiaID string) (*models.KeluargaLengkap, error)
}

type keluargaRepository struct {
	db *sqlx.DB
}

func newKeluargaRepository(db *sqlx.DB) KeluargaRepository {
	return &keluargaRepository{db: db}
}

func (r *keluargaRepository) GetDashboard(ctx context.Context, lansiaID string) (*models.KeluargaDashboard, error) {
	var lansia models.Lansia
	if err := r.db.GetContext(ctx, &lansia, `SELECT * FROM lansia WHERE id = $1`, lansiaID); err != nil {
		return nil, err
	}

	var catatanHariIni []models.CatatanShiftPublik
	r.db.SelectContext(ctx, &catatanHariIni, `
		SELECT cs.id, cs.narasi_shift, cs.suasana_hati, cs.nafsu_makan,
		       cs.aktivitas, cs.nama_shift, u.name AS nama_pengurus, cs.created_at
		FROM catatan_shift cs
		JOIN users u ON u.id = cs.pengurus_id
		WHERE cs.lansia_id = $1 AND DATE(cs.created_at) = CURRENT_DATE
		ORDER BY cs.created_at DESC
	`, lansiaID)

	var vitalTerbaru models.VitalPublik
	r.db.GetContext(ctx, &vitalTerbaru, `
		SELECT tekanan_darah, fase_label, fase_warna, suhu_tubuh, gula_darah, created_at
		FROM pemeriksaan_kesehatan WHERE lansia_id = $1 ORDER BY created_at DESC LIMIT 1
	`, lansiaID)

	var obatPagi []models.ObatStatus
	r.db.SelectContext(ctx, &obatPagi, `
		SELECT o.nama_obat, j.dosis,
		       (lp.id IS NOT NULL) AS sudah_diberikan,
		       COALESCE(TO_CHAR(lp.diberikan_pada, 'HH24:MI'), '') AS waktu_diberikan
		FROM jadwal_obat j
		JOIN obat o ON o.id = j.obat_id
		LEFT JOIN log_pemberian_obat lp
			ON lp.jadwal_obat_id = j.id AND DATE(lp.diberikan_pada) = CURRENT_DATE
		WHERE o.lansia_id = $1 AND j.sesi = 'pagi'
		ORDER BY o.nama_obat
	`, lansiaID)

	var riwayatKunjungan []models.KunjunganRingkas
	r.db.SelectContext(ctx, &riwayatKunjungan, `
		SELECT nama_keluarga, hubungan_keluarga, tanggal_kunjungan, durasi_menit, catatan
		FROM kunjungan_keluarga WHERE lansia_id = $1 ORDER BY tanggal_kunjungan DESC LIMIT 3
	`, lansiaID)

	var kondisiFisik models.KondisiFisikRingkas
	r.db.GetContext(ctx, &kondisiFisik, `
		SELECT jenis_kondisi, lokasi_tubuh, risiko_label, tingkat_risiko, foto_url, created_at
		FROM galeri_fisik WHERE lansia_id = $1 ORDER BY created_at DESC LIMIT 1
	`, lansiaID)

	return &models.KeluargaDashboard{
		Lansia:              lansia,
		CatatanHariIni:      catatanHariIni,
		VitalTerbaru:        vitalTerbaru,
		ObatPagi:            obatPagi,
		RiwayatKunjungan:    riwayatKunjungan,
		KondisiFisikTerbaru: kondisiFisik,
	}, nil
}

func (r *keluargaRepository) GetKeluarga(ctx context.Context, lansiaID string) (*models.KeluargaLengkap, error) {
	var lansia models.Lansia
	if err := r.db.GetContext(ctx, &lansia, `SELECT * FROM lansia WHERE id = $1`, lansiaID); err != nil {
		return nil, err
	}

	var catatanShift []models.CatatanShift
	r.db.SelectContext(ctx, &catatanShift, `
		SELECT cs.*, u.name AS nama_pengurus FROM catatan_shift cs
		JOIN users u ON u.id = cs.pengurus_id
		WHERE cs.lansia_id = $1 AND cs.created_at >= NOW() - INTERVAL '30 days'
		ORDER BY cs.created_at DESC
	`, lansiaID)

	var riwayatPemeriksaan []models.PemeriksaanKesehatan
	r.db.SelectContext(ctx, &riwayatPemeriksaan, `
		SELECT * FROM pemeriksaan_kesehatan WHERE lansia_id = $1 ORDER BY created_at DESC LIMIT 10
	`, lansiaID)

	var obatAktif []models.Obat
	r.db.SelectContext(ctx, &obatAktif, `SELECT * FROM obat WHERE lansia_id = $1 AND is_aktif = true`, lansiaID)

	var totalObat, sudahDiberikan int
	r.db.GetContext(ctx, &totalObat, `
		SELECT COUNT(*) FROM jadwal_obat j JOIN obat o ON o.id = j.obat_id WHERE o.lansia_id = $1
	`, lansiaID)
	r.db.GetContext(ctx, &sudahDiberikan, `
		SELECT COUNT(DISTINCT lp.jadwal_obat_id) FROM log_pemberian_obat lp
		JOIN jadwal_obat j ON j.id = lp.jadwal_obat_id
		JOIN obat o ON o.id = j.obat_id
		WHERE o.lansia_id = $1 AND DATE(lp.diberikan_pada) = CURRENT_DATE
	`, lansiaID)

	var riwayatGaleri []models.GaleriFisik
	r.db.SelectContext(ctx, &riwayatGaleri, `
		SELECT * FROM galeri_fisik WHERE lansia_id = $1 ORDER BY created_at DESC LIMIT 10
	`, lansiaID)

	var riwayatKunjungan []models.KunjunganKeluarga
	r.db.SelectContext(ctx, &riwayatKunjungan, `
		SELECT * FROM kunjungan_keluarga WHERE lansia_id = $1 ORDER BY tanggal_kunjungan DESC
	`, lansiaID)

	return &models.KeluargaLengkap{
		Lansia:             lansia,
		CatatanShift:       catatanShift,
		RiwayatPemeriksaan: riwayatPemeriksaan,
		ObatAktif:          obatAktif,
		ObatHariIni: models.ObatHariIniSummary{
			Total:          totalObat,
			SudahDiberikan: sudahDiberikan,
			BelumDiberikan: totalObat - sudahDiberikan,
		},
		RiwayatKondisiFisik: riwayatGaleri,
		RiwayatKunjungan:    riwayatKunjungan,
	}, nil
}
