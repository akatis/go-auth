package main

import (
	"github.com/akatis/go-auth/test/authTest"
	"github.com/akatis/go-auth/test/handler"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {

	app := fiber.New()
	api := app.Group("/api")

	a := authTest.GetConf()

	api.Use(a.Middleware)

	app.Get("/get", handler.GetToken)
	api.Get("/test", handler.Test)

	err := app.Listen(":7777")
	if err != nil {
		log.Fatal(err)
	}
}
