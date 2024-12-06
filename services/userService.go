package services

import (
	"ObserverKVS/models"
	"ObserverKVS/repositories"

	"github.com/gofiber/fiber/v2"
)

func HandleCreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	userRepo, err := repositories.NewUserRepository()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	apikey, err := userRepo.CreateUser(user)
	if err != nil {
		return c.SendStatus(fiber.StatusConflict)

	}
	return c.JSON(apikey)
}

func HandlerGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userRepo, err := repositories.NewUserRepository()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user, err := userRepo.GetUserByID(id)
	if err != nil {
		if err.Error() == "notfound" {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if user == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(user)
}

func HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userRepo, err := repositories.NewUserRepository()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = userRepo.DeleteUser(id)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
