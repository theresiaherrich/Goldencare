func RequireRole(roles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Ambil role dari locals (sudah di-set oleh RequireAuth)
        userRole, ok := c.Locals("role").(string)
        if !ok || userRole == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "unauthorized",
            })
        }

        // Superadmin bypass semua role check
        if userRole == "superadmin" {
            return c.Next()
        }

        // Cek role yang diizinkan
        for _, role := range roles {
            if userRole == role {
                return c.Next()
            }
        }

        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "forbidden: insufficient role",
        })
    }
}