package handler

import (
	"encoding/json"
	"fmt"
	"github.com/akatis/go-auth/test/authTest"
	"github.com/gofiber/fiber/v2"
)

func Test(ctx *fiber.Ctx) error {
	a := authTest.GetAuth()
	fmt.Println("api")
	fmt.Println(ctx.Route().Path)

	uuid, _ := a.GetUUID(ctx)

	data, _ := json.Marshal(uuid)
	ctx.Response().SetStatusCode(200)
	ctx.Response().Header.Add("Content-Type", "application/json")
	ctx.Write(data)
	err := a.DeleteFromRedis(a.Payload)

	if err != nil {
		return err
	}
	return nil
}
