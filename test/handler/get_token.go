package handler

import (
	"encoding/json"
	"github.com/akatis/go-auth/test/authTest"
	"github.com/gofiber/fiber/v2"
)

func GetToken(ctx *fiber.Ctx) error {
	a := authTest.GetAuth()

	token := a.CreateAccessToken("uuid3", []int{0})

	_ = a.AddToRedis("uuid3", "user agent")

	data, _ := json.Marshal(token)
	ctx.Response().SetStatusCode(200)
	ctx.Response().Header.Add("Content-Type", "application/json")
	ctx.Write(data)
	return nil
}
