package repository

import "github.com/jmoiron/sqlx"

type Repository interface {
	User() UserRepository
	KodeUndangan() KodeUndanganRepository
	Panti() PantiRepository
	Lansia() LansiaRepository
	Superadmin() SuperadminRepository
	Galeri() GaleriRepository
	CatatanShift() CatatanShiftRepository
	Pemeriksaan() PemeriksaanRepository
	Kunjungan() KunjunganRepository
	Obat() ObatRepository
	Pengurus() PengurusRepository
	Keluarga() KeluargaRepository
	Close() error
}

type postgresRepository struct {
	db               *sqlx.DB
	userRepo         UserRepository
	kodeUndanganRepo KodeUndanganRepository
	pantiRepo        PantiRepository
	lansiaRepo       LansiaRepository
	superadminRepo   SuperadminRepository
	galeriRepo       GaleriRepository
	catatanShiftRepo CatatanShiftRepository
	pemeriksaanRepo  PemeriksaanRepository
	kunjunganRepo    KunjunganRepository
	obatRepo         ObatRepository
	pengurusRepo     PengurusRepository
	keluargaRepo     KeluargaRepository
}

func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{
		db:               db,
		userRepo:         newUserRepository(db),
		kodeUndanganRepo: newKodeUndanganRepository(db),
		pantiRepo:        newPantiRepository(db),
		lansiaRepo:       newLansiaRepository(db),
		superadminRepo:   newSuperadminRepository(db),
		galeriRepo:       newGaleriRepository(db),
		catatanShiftRepo: newCatatanShiftRepository(db),
		pemeriksaanRepo:  newPemeriksaanRepository(db),
		kunjunganRepo:    newKunjunganRepository(db),
		obatRepo:         newObatRepository(db),
		pengurusRepo:     newPengurusRepository(db),
		keluargaRepo:     newKeluargaRepository(db),
	}
}

func (r *postgresRepository) User() UserRepository                 { return r.userRepo }
func (r *postgresRepository) KodeUndangan() KodeUndanganRepository { return r.kodeUndanganRepo }
func (r *postgresRepository) Panti() PantiRepository               { return r.pantiRepo }
func (r *postgresRepository) Lansia() LansiaRepository             { return r.lansiaRepo }
func (r *postgresRepository) Superadmin() SuperadminRepository     { return r.superadminRepo }
func (r *postgresRepository) Galeri() GaleriRepository             { return r.galeriRepo }
func (r *postgresRepository) CatatanShift() CatatanShiftRepository { return r.catatanShiftRepo }
func (r *postgresRepository) Pemeriksaan() PemeriksaanRepository   { return r.pemeriksaanRepo }
func (r *postgresRepository) Kunjungan() KunjunganRepository       { return r.kunjunganRepo }
func (r *postgresRepository) Obat() ObatRepository                 { return r.obatRepo }
func (r *postgresRepository) Pengurus() PengurusRepository         { return r.pengurusRepo }
func (r *postgresRepository) Keluarga() KeluargaRepository         { return r.keluargaRepo }
func (r *postgresRepository) Close() error                         { return r.db.Close() }
