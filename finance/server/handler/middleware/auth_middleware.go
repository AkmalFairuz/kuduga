package middleware

import (
	"github.com/akmalfairuz/finance/server/handler/responses"
	"github.com/akmalfairuz/finance/server/handler/xhelper"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/gofiber/fiber/v2"
)

func onAuthFailed(ctx *fiber.Ctx) error {
	ctx.Response().Header.Add("X-Require-Auth", "true")
	return ctx.Status(fiber.StatusUnauthorized).JSON(responses.UnauthorizedError)
}

func AuthMiddleware(authService *service.AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorization := string(ctx.Request().Header.Peek("Authorization"))

		tok := xhelper.GetToken(authService, authorization)
		if tok == nil {
			return onAuthFailed(ctx)
		}

		ctx.Locals("token", *tok)

		return ctx.Next()
	}
}
