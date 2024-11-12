package middlewareLib

import (
	"latihan-portal-news/config"
	authLib "latihan-portal-news/lib/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Middleware interface {
	CheckToken() fiber.Handler
}

type Options struct {
	authJwt authLib.Jwt
}

// CheckToken implements Middleware.
func (o *Options) CheckToken() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHandler := c.Get("Authorization")
		if authHandler == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is required",
			})
		}

		tokenString := strings.Split(authHandler, "Bearer ")[1]
		claims, err := o.authJwt.VerifyAccessToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		c.Locals("user", claims)

		return c.Next()
	}
}

func NewMiddleware(cfg *config.Config) Middleware {
	opt := new(Options)
	opt.authJwt = authLib.NewJwt(cfg)

	return opt
}
