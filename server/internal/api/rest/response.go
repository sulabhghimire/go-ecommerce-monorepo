package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorMessage(ctx *fiber.Ctx, statusCode int, err error) error {
	return ctx.Status(statusCode).JSON(fiber.Map{"message": err.Error()})
}

func NotFoundError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
}

func InternalError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
}

func BadRequest(ctx *fiber.Ctx, msg string) error {
	return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"message": msg})
}

func NotAuhtorizedError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
}

func SuccessResponse(ctx *fiber.Ctx, statusCode int, msg string, data interface{}) error {
	return ctx.Status(statusCode).JSON(&fiber.Map{
		"message": msg,
		"data":    data,
	})
}
