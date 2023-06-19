package handler

import (
	"encoding/json"
	"github.com/akatis/go-auth/test/authTest"
	"github.com/gofiber/fiber/v2"
)

func Test(ctx *fiber.Ctx) error {
	a := authTest.GetAuth()

	uuid, _ := a.GetUUID(ctx)

	data, _ := json.Marshal(uuid)
	ctx.Response().SetStatusCode(200)
	ctx.Response().Header.Add("Content-Type", "application/json")
	ctx.Write(data)
	a.DeleteFromRedis(a.Payload)
	return nil
}
