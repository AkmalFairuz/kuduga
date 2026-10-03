package ctxhelper

import (
	"github.com/akmalfairuz/finance/module/password"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/gofiber/fiber/v2"
)

func GetToken(ctx *fiber.Ctx) model.Token {
	return ctx.Locals("token").(model.Token)
}

func GetUserID(ctx *fiber.Ctx) int64 {
	return GetToken(ctx).UserID
}

func IsAdmin(ctx *fiber.Ctx) bool {
	return GetToken(ctx).Role == model.RoleAdmin
}

func ValidatePin(pinHash *string, pin string) bool {
	if pinHash == nil {
		return pin == ""
	}
	err := password.Verify([]byte(*pinHash), []byte(pin))
	return err == nil
}
