package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/theresiaherrich/Goldencare/internal/config"
	"github.com/theresiaherrich/Goldencare/pkg/utils"
)

type Claims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	PantiID string `json:"panti_id"`
	jwt.RegisteredClaims
}

func Logger() fiber.Handler {
	return fiberLogger.New()
}

func RequireAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.Unauthorized(c, "Token tidak ditemukan")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.Unauthorized(c, "Format token tidak valid")
		}

		tokenStr := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return utils.Unauthorized(c, "Token tidak valid atau sudah kedaluwarsa")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("panti_id", claims.PantiID)

		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok || userRole == "" {
			return utils.Unauthorized(c, "Role tidak ditemukan")
		}

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return utils.Forbidden(c, "Anda tidak memiliki akses ke fitur ini")
	}
}

func GetUserID(c *fiber.Ctx) string {
	id, _ := c.Locals("user_id").(string)
	return id
}

func GetRole(c *fiber.Ctx) string {
	role, _ := c.Locals("role").(string)
	return role
}

func GetPantiID(c *fiber.Ctx) string {
	pantiID := c.Locals("panti_id")
	if pantiID == nil {
		return ""
	}
	return pantiID.(string)
}
