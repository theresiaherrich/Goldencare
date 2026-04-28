package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/theresiaherrich/Goldencare/internal/middleware"
    "github.com/theresiaherrich/Goldencare/internal/models"
    "github.com/theresiaherrich/Goldencare/internal/services"
    "github.com/theresiaherrich/Goldencare/pkg/utils"
)

type AuthHandler struct {
    authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req models.RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.BadRequest(c, "Request body tidak valid")
    }

    user, err := h.authService.Register(c.Context(), &req)
    if err != nil {
        return utils.BadRequest(c, err.Error())
    }

    return utils.Created(c, "Registrasi berhasil", models.ToUserResponse(user))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req models.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.BadRequest(c, "Request body tidak valid")
    }

    user, token, err := h.authService.Login(c.Context(), &req)
    if err != nil {
        return utils.Unauthorized(c, err.Error())
    }

    return utils.OK(c, "Login berhasil", models.LoginResponse{
        User:  models.ToUserResponse(user),
        Token: token,
    })
}

func (h *AuthHandler) SuperadminLogin(c *fiber.Ctx) error {
    var req models.SuperadminLoginRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.BadRequest(c, "Request body tidak valid")
    }

    superadmin, token, err := h.authService.SuperadminLogin(c.Context(), &req)
    if err != nil {
        return utils.Unauthorized(c, err.Error())
    }

    return utils.OK(c, "Superadmin login berhasil", models.SuperadminLoginResponse{
        User:  models.ToSuperadminResponse(superadmin),
        Token: token,
    })
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
    userID := middleware.GetUserID(c)
    if userID == "" {
        return utils.Unauthorized(c, "User ID tidak ditemukan")
    }

    user, err := h.authService.Me(c.Context(), userID)
    if err != nil {
        return utils.NotFound(c, "User tidak ditemukan")
    }

    return utils.OK(c, "Data user", models.ToUserResponse(user))
}

func (h *AuthHandler) GenerateKode(c *fiber.Ctx) error {
    userID := middleware.GetUserID(c)
    pantiID := middleware.GetPantiID(c)

    if userID == "" || pantiID == "" {
        return utils.BadRequest(c, "User ID atau Panti ID tidak ditemukan")
    }

    var req models.GenerateKodeRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.BadRequest(c, "Request body tidak valid")
    }

    kode, err := h.authService.GenerateKode(c.Context(), userID, pantiID, &req)
    if err != nil {
        return utils.BadRequest(c, err.Error())
    }

    return utils.Created(c, "Kode undangan berhasil dibuat", models.ToKodeResponse(kode))
}

func (h *AuthHandler) SuperadminGenerateKode(c *fiber.Ctx) error {
    var req models.GenerateKodePengelolaRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.BadRequest(c, "Request body tidak valid")
    }

    kode, err := h.authService.SuperadminGenerateKode(c.Context(), &req)
    if err != nil {
        return utils.BadRequest(c, err.Error())
    }

    return utils.Created(c, "Kode undangan berhasil dibuat", models.ToKodeResponse(kode))
}

func (h *AuthHandler) ListKode(c *fiber.Ctx) error {
    pantiID := middleware.GetPantiID(c)
    if pantiID == "" {
        return utils.BadRequest(c, "Panti ID tidak ditemukan")
    }

    kodes, err := h.authService.ListKode(c.Context(), pantiID)
    if err != nil {
        return utils.InternalError(c, "Gagal mengambil daftar kode")
    }

    return utils.OK(c, "Daftar kode undangan", kodes)
}

func (h *AuthHandler) NonaktifkanKode(c *fiber.Ctx) error {
    kodeID := c.Params("kode")
    if kodeID == "" {
        return utils.BadRequest(c, "Kode ID tidak ditemukan")
    }

    if _, err := uuid.Parse(kodeID); err != nil {
        return utils.BadRequest(c, "Format kode ID tidak valid")
    }

    if err := h.authService.NonaktifkanKode(c.Context(), kodeID); err != nil {
        return utils.BadRequest(c, err.Error())
    }

    return utils.OK(c, "Kode undangan berhasil dinonaktifkan", nil)
}