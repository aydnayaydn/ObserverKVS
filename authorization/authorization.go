package authorization

import (
	"ObserverKVS/repositories"

	"github.com/gofiber/fiber/v2"
)

func AdminAuth(fn fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRepo, err := repositories.NewUserRepository()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		apiKey := c.Get("api-key")
		if apiKey == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		user, err := userRepo.GetUserByApiKey(apiKey)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if user == nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		if user.Role != "admin" {
			return c.SendStatus(fiber.StatusForbidden)
		}

		return fn(c)
	}
}

func StackholdersAuth(fn fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRepo, err := repositories.NewUserRepository()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		apiKey := c.Get("api-key")
		if apiKey == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		user, err := userRepo.GetUserByApiKey(apiKey)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if user == nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		if user.Role == "stackholder" || user.Role == "admin" {
			return fn(c)
		}

		return c.SendStatus(fiber.StatusUnauthorized)
	}
}
