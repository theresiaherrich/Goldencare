package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `db:"id"           json:"id"`
	Name        string     `db:"name"         json:"name"`
	Email       string     `db:"email"        json:"email"`
	Password    string     `db:"password"     json:"-"`
	Role        string     `db:"role"         json:"role"`
	IsVerified  bool       `db:"is_verified"  json:"is_verified"`
	PantiID     *uuid.UUID `db:"panti_id"     json:"panti_id"`
	KodeDipakai string     `db:"kode_dipakai" json:"kode_dipakai,omitempty"`
	CreatedAt   time.Time  `db:"created_at"   json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"   json:"updated_at"`
}

type Superadmin struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	Email     string    `db:"email"      json:"email"`
	Password  string    `db:"password"   json:"-"`
	Name      string    `db:"name"       json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type KodeUndangan struct {
	ID           uuid.UUID  `db:"id"            json:"id"`
	Kode         string     `db:"kode"          json:"kode"`
	UntukRole    string     `db:"untuk_role"    json:"untuk_role"`
	DibuatOleh   *uuid.UUID `db:"dibuat_oleh"   json:"dibuat_oleh"`
	PantiID      *uuid.UUID `db:"panti_id"      json:"panti_id"`
	Tipe         string     `db:"tipe"          json:"tipe"`
	DipakaiCount int        `db:"dipakai_count" json:"dipakai_count"`
	MaksPakai    *int       `db:"maks_pakai"    json:"maks_pakai"`
	ExpiredAt    *time.Time `db:"expired_at"    json:"expired_at"`
	IsAktif      bool       `db:"is_aktif"      json:"is_aktif"`
	Catatan      string     `db:"catatan"       json:"catatan"`
	CreatedAt    time.Time  `db:"created_at"    json:"created_at"`
}

type RegisterRequest struct {
	Name         string `json:"name"          validate:"required,min=2"`
	Email        string `json:"email"         validate:"required,email"`
	Password     string `json:"password"      validate:"required,min=6"`
	Role         string `json:"role"          validate:"required,oneof=pengelola pengurus keluarga"`
	KodeUndangan string `json:"kode_undangan" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SuperadminLoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type JoinPantiRequest struct {
	KodeUndangan string `json:"kode_undangan" validate:"required"`
}

type GenerateKodeRequest struct {
	UntukRole string `json:"untuk_role" validate:"required,oneof=pengurus keluarga"`
	Tipe      string `json:"tipe"       validate:"required,oneof=single_use multi_use"`
	MaksPakai *int   `json:"maks_pakai"`
	ExpiredAt string `json:"expired_at"`
	Catatan   string `json:"catatan"`
}

type GenerateKodePengelolaRequest struct {
	Tipe      string `json:"tipe"      validate:"required,oneof=single_use multi_use"`
	MaksPakai *int   `json:"maks_pakai"`
	ExpiredAt string `json:"expired_at"`
	Catatan   string `json:"catatan"`
}

type UserResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	PantiID    *uuid.UUID `json:"panti_id"`
	IsVerified bool       `json:"is_verified"`
}

type SuperadminResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type SuperadminLoginResponse struct {
	User  SuperadminResponse `json:"user"`
	Token string             `json:"token"`
}

type KodeResponse struct {
	ID   uuid.UUID `json:"id"`
	Kode string    `json:"kode"`
	Tipe string    `json:"tipe"`
}

func ToUserResponse(u *User) UserResponse {
	return UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Role:       u.Role,
		PantiID:    u.PantiID,
		IsVerified: u.IsVerified,
	}
}

func ToSuperadminResponse(s *Superadmin) SuperadminResponse {
	return SuperadminResponse{
		ID:    s.ID,
		Name:  s.Name,
		Email: s.Email,
	}
}

func ToKodeResponse(k *KodeUndangan) KodeResponse {
	return KodeResponse{
		ID:   k.ID,
		Kode: k.Kode,
		Tipe: k.Tipe,
	}
}