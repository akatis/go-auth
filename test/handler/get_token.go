package handler

import (
	"encoding/json"
	"github.com/akatis/go-auth/v3/test/authTest"
	"github.com/gofiber/fiber/v2"
)

func GetToken(ctx *fiber.Ctx) error {
	a := authTest.GetAuth()

	s := 44
	token := a.CreateAccessToken("uuid3", []int{99}, nil, &s)

	data, _ := json.Marshal(token)
	ctx.Response().SetStatusCode(200)
	ctx.Response().Header.Add("Content-Type", "application/json")
	ctx.Write(data)
	return nil
}
