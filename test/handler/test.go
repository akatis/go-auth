package handler

import (
	"encoding/json"
	"github.com/akatis/go-auth/test/authTest"
	"github.com/gofiber/fiber/v2"
)

func Test(ctx *fiber.Ctx) error {
	a := authTest.GetConf()

	data, _ := json.Marshal("It works too‼")
	ctx.Response().SetStatusCode(200)
	ctx.Response().Header.Add("Content-Type", "application/json")
	ctx.Write(data)
	a.DeleteFromRedis(a.Payload)
	return nil
}
