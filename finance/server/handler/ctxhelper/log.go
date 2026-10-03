package ctxhelper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Log(ctx *fiber.Ctx) logrus.FieldLogger {
	return ctx.Locals("log").(logrus.FieldLogger)
}
