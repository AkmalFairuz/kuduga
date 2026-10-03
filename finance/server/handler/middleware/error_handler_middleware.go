package middleware

import (
	"errors"
	"github.com/akmalfairuz/finance/module/random"
	"github.com/akmalfairuz/finance/server/handler/ctxhelper"
	"github.com/akmalfairuz/finance/server/handler/responses"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandlerMiddleware() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		var e errortype.VisibleError
		if errors.As(err, &e) {
			return e.Print(ctx)
		}
		errorId := random.StringUpper(8)
		ctxhelper.Log(ctx).Errorf("error %s: %+v", errorId, err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.M("Terjadi kesalahan di sisi server\nError ID: " + errorId))
	}
}
