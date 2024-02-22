package main

import (
	"github.com/akatis/go-auth/test/authTest"
	"github.com/akatis/go-auth/test/handler"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	a := authTest.GetAuth()

	app := fiber.New()
	app.Get("/get", handler.GetToken)

	api := app.Group("/api")
	api.Use(a.Middleware)

	api.Get("/testo", handler.Test)
	api.Get("/test/:id/user/:name", handler.Test)

	err := app.Listen(":7777")
	if err != nil {
		log.Fatal(err)
	}
}
