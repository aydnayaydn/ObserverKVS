package services

import (
	"errors"

	"ObserverKVS/models"
	"ObserverKVS/repositories"

	"github.com/gofiber/fiber/v2"
)

func HandleSaveKeyValuePair(c *fiber.Ctx) error {
	pair := new(models.KeyValuePair)
	if err := c.BodyParser(pair); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	stackholderName, err := getStackholderName(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	keyValuePairRepo, err := repositories.NewKeyValuePairRepository(stackholderName)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	id, err := keyValuePairRepo.Save(pair)
	if err != nil {
		return c.SendStatus(fiber.StatusConflict)

	}

	return c.JSON(id)
}

func HandlerGetKeyValuePair(c *fiber.Ctx) error {
	key := c.Params("key")

	stackholderName, err := getStackholderName(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	keyValuePairRepo, err := repositories.NewKeyValuePairRepository(stackholderName)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	value, err := keyValuePairRepo.Get(key)
	if err != nil {
		if err.Error() == "notfound" {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if value == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.JSON(value)
}

func HandleDeleteKeyValuePair(c *fiber.Ctx) error {
	key := c.Params("key")

	stackholderName, err := getStackholderName(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	keyValuePairRepo, err := repositories.NewKeyValuePairRepository(stackholderName)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = keyValuePairRepo.Delete(key)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func HandleUpdateKeyValuePair(c *fiber.Ctx) error {
	pair := new(models.KeyValuePair)
	if err := c.BodyParser(pair); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	stackholderName, err := getStackholderName(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	keyValuePairRepo, err := repositories.NewKeyValuePairRepository(stackholderName)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	isSuccess, err := keyValuePairRepo.Update(pair)
	if err != nil {
		return c.SendStatus(fiber.StatusConflict)
	}

	if isSuccess {
		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.SendStatus(fiber.StatusNotFound)
	}
}

func HandlerGetKeyValuePairById(c *fiber.Ctx) error {
	id := c.Params("id")

	stackholderName, err := getStackholderName(c)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	keyValuePairRepo, err := repositories.NewKeyValuePairRepository(stackholderName)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	value, err := keyValuePairRepo.GetById(id)
	if err != nil {
		if err.Error() == "notfound" {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if value == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.JSON(value)
}

func getStackholderName(c *fiber.Ctx) (string, error) {
	userRepo, err := repositories.NewUserRepository()
	if err != nil {
		return "", errors.New("Internal Server Error")
	}

	apikey := c.Get("api-key")
	if apikey == "" {
		return "", errors.New("Bad Request")
	}

	user, err := userRepo.GetUserByApiKey(apikey)
	if err != nil {
		return "", errors.New("Internal Server Error")
	}

	if user == nil {
		return "", errors.New("Not Found")
	}

	return user.Username, nil
}
