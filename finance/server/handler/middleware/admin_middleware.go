package middleware

import (
	"github.com/akmalfairuz/finance/server/handler/ctxhelper"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := ctxhelper.GetToken(ctx)
		if token.Role != model.RoleAdmin {
			return fiber.ErrForbidden
		}
		return ctx.Next()
	}
}
