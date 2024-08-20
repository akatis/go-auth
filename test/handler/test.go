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

	//uuid, _ := a.GetUUID(ctx)
	shopId, err := a.GetShopID(ctx)

	if err != nil {
		data, _ := json.Marshal(err.Error())
		ctx.Response().SetStatusCode(200)
		ctx.Response().Header.Add("Content-Type", "application/json")
		ctx.Write(data)
		return nil
	}
	data, _ := json.Marshal(shopId)
	ctx.Response().SetStatusCode(200)
	ctx.Response().Header.Add("Content-Type", "application/json")
	ctx.Write(data)

	return nil
}
