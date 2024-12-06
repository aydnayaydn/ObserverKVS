package main

import (
	"log"

	"ObserverKVS/authorization"
	"ObserverKVS/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/connection", ConnectionHandler())

	app.Post("/user/", authorization.AdminAuth(services.HandleCreateUser))
	app.Get("/user/:id", authorization.AdminAuth(services.HandlerGetUser))
	app.Delete("/user/:id", authorization.AdminAuth(services.HandleDeleteUser))

	app.Post("/keyvaluepair/", authorization.StackholdersAuth(services.HandleSaveKeyValuePair))
	app.Get("/keyvaluepair/:key", authorization.StackholdersAuth(services.HandlerGetKeyValuePair))
	app.Get("/keyvaluepair/id/:id", authorization.StackholdersAuth(services.HandlerGetKeyValuePairById))
	app.Put("/keyvaluepair/", authorization.StackholdersAuth(services.HandleUpdateKeyValuePair))
	app.Delete("/keyvaluepair/:key", authorization.StackholdersAuth(services.HandleDeleteKeyValuePair))

	log.Fatal(app.Listen(":3000"))
}

func ConnectionHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Connection is successful")
	}
}
