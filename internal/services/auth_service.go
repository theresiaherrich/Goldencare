package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/theresiaherrich/Goldencare/internal/config"
	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaims struct {
	UserID  string
	Email   string
	Role    string
	PantiID string
}

type jwtTokenClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	PantiID string `json:"panti_id"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.User, string, error)
	SuperadminLogin(ctx context.Context, req *models.SuperadminLoginRequest) (*models.Superadmin, string, error)
	ValidateToken(ctx context.Context, token string) (*JWTClaims, error)
	GenerateKode(ctx context.Context, userID string, pantiID string, req *models.GenerateKodeRequest) (*models.KodeUndangan, error)
	SuperadminGenerateKode(ctx context.Context, req *models.GenerateKodePengelolaRequest) (*models.KodeUndangan, error)
	ListKode(ctx context.Context, pantiID string) ([]models.KodeUndangan, error)
	NonaktifkanKode(ctx context.Context, kodeID string) error
	Me(ctx context.Context, userID string) (*models.User, error)
}

type authService struct {
	repo repository.Repository
	cfg  *config.Config
}

func NewAuthService(repo repository.Repository, cfg *config.Config) AuthService {
	return &authService{repo: repo, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	validRoles := map[string]bool{"pengelola": true, "pengurus": true, "keluarga": true}
	if !validRoles[req.Role] {
		return nil, fmt.Errorf("role tidak valid")
	}

	existingUser, _ := s.repo.User().GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email sudah terdaftar")
	}

	kode, err := s.repo.KodeUndangan().GetByKode(ctx, req.KodeUndangan)
	if err != nil {
		return nil, fmt.Errorf("kode undangan tidak valid")
	}
	if !kode.IsAktif {
		return nil, fmt.Errorf("kode undangan tidak aktif")
	}
	if kode.UntukRole != req.Role && kode.UntukRole != "all" {
		return nil, fmt.Errorf("kode undangan tidak sesuai dengan role")
	}
	if kode.MaksPakai != nil && kode.DipakaiCount >= *kode.MaksPakai {
		return nil, fmt.Errorf("kode undangan sudah mencapai batas penggunaan")
	}
	if kode.ExpiredAt != nil && kode.ExpiredAt.Before(time.Now()) {
		return nil, fmt.Errorf("kode undangan sudah kadaluarsa")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal hash password: %w", err)
	}

	user := &models.User{
		ID:          uuid.New(),
		Name:        req.Name,
		Email:       strings.ToLower(req.Email),
		Password:    string(hashedPassword),
		Role:        req.Role,
		IsVerified:  false,
		KodeDipakai: req.KodeUndangan,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.Role == "pengelola" {
		panti := &models.Panti{
			ID:          uuid.New(),
			Nama:        req.Name,
			PengelolaID: user.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := s.repo.Panti().Create(ctx, panti); err != nil {
			return nil, fmt.Errorf("gagal membuat panti: %w", err)
		}
		user.PantiID = &panti.ID
	}

	if err := s.repo.User().Create(ctx, user); err != nil {
		return nil, fmt.Errorf("gagal membuat user: %w", err)
	}

	kode.DipakaiCount++
	if err := s.repo.KodeUndangan().Update(ctx, kode); err != nil {
		return nil, fmt.Errorf("gagal update kode: %w", err)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.User, string, error) {
	user, err := s.repo.User().GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("email atau password salah")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", fmt.Errorf("email atau password salah")
	}
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("gagal generate token: %w", err)
	}
	return user, token, nil
}

func (s *authService) SuperadminLogin(ctx context.Context, req *models.SuperadminLoginRequest) (*models.Superadmin, string, error) {
	superadmin, err := s.repo.Superadmin().GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("email atau password salah")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(superadmin.Password), []byte(req.Password)); err != nil {
		return nil, "", fmt.Errorf("email atau password salah")
	}
	token, err := s.generateSuperadminToken(superadmin)
	if err != nil {
		return nil, "", fmt.Errorf("gagal generate token: %w", err)
	}
	return superadmin, token, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
	claims := &jwtTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token tidak valid")
	}
	return &JWTClaims{
		UserID:  claims.UserID,
		Email:   claims.Email,
		Role:    claims.Role,
		PantiID: claims.PantiID,
	}, nil
}

func (s *authService) GenerateKode(ctx context.Context, userID string, pantiID string, req *models.GenerateKodeRequest) (*models.KodeUndangan, error) {
	userUUID := uuid.MustParse(userID)
	pantiUUID := uuid.MustParse(pantiID)

	kode := &models.KodeUndangan{
		ID:           uuid.New(),
		Kode:         generateRandomKode(),
		UntukRole:    req.UntukRole,
		DibuatOleh:   &userUUID,
		PantiID:      &pantiUUID,
		Tipe:         req.Tipe,
		DipakaiCount: 0,
		MaksPakai:    req.MaksPakai,
		IsAktif:      true,
		Catatan:      req.Catatan,
		CreatedAt:    time.Now(),
	}
	if req.ExpiredAt != "" {
		if t, err := time.Parse("2006-01-02", req.ExpiredAt); err == nil {
			kode.ExpiredAt = &t
		}
	}
	if err := s.repo.KodeUndangan().Create(ctx, kode); err != nil {
		return nil, fmt.Errorf("gagal membuat kode: %w", err)
	}
	return kode, nil
}

func (s *authService) SuperadminGenerateKode(ctx context.Context, req *models.GenerateKodePengelolaRequest) (*models.KodeUndangan, error) {
	kode := &models.KodeUndangan{
		ID:           uuid.New(),
		Kode:         generateRandomKode(),
		UntukRole:    "pengelola",
		Tipe:         req.Tipe,
		DipakaiCount: 0,
		MaksPakai:    req.MaksPakai,
		IsAktif:      true,
		Catatan:      req.Catatan,
		CreatedAt:    time.Now(),
	}
	if req.ExpiredAt != "" {
		if t, err := time.Parse("2006-01-02", req.ExpiredAt); err == nil {
			kode.ExpiredAt = &t
		}
	}
	if err := s.repo.KodeUndangan().Create(ctx, kode); err != nil {
		return nil, fmt.Errorf("gagal membuat kode: %w", err)
	}
	return kode, nil
}

func (s *authService) ListKode(ctx context.Context, pantiID string) ([]models.KodeUndangan, error) {
	return s.repo.KodeUndangan().GetByPantiID(ctx, pantiID)
}

func (s *authService) NonaktifkanKode(ctx context.Context, kodeID string) error {
	kode, err := s.repo.KodeUndangan().GetByID(ctx, kodeID)
	if err != nil {
		return fmt.Errorf("kode tidak ditemukan")
	}
	kode.IsAktif = false
	return s.repo.KodeUndangan().Update(ctx, kode)
}

func (s *authService) Me(ctx context.Context, userID string) (*models.User, error) {
	return s.repo.User().GetByID(ctx, userID)
}

func (s *authService) generateToken(user *models.User) (string, error) {
	pantiID := ""
	if user.PantiID != nil {
		pantiID = user.PantiID.String()
	}
	claims := jwtTokenClaims{
		UserID:  user.ID.String(),
		Email:   user.Email,
		Role:    user.Role,
		PantiID: pantiID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.cfg.JWTExpireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *authService) generateSuperadminToken(superadmin *models.Superadmin) (string, error) {
	claims := jwtTokenClaims{
		UserID:  superadmin.ID.String(),
		Email:   superadmin.Email,
		Role:    "superadmin",
		PantiID: "",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.cfg.JWTExpireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func generateRandomKode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}