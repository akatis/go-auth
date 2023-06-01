package auth

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (response *Response) HttpResponse(ctx *fiber.Ctx, status int) error {
	ctx.Status(status)

	ctx.Response().Header.Add("Content-Type", "application/json")
	w, _ := json.Marshal(response)
	ctx.Write(w)

	return nil
}
