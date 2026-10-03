package ctxhelper

import (
	"github.com/akmalfairuz/finance/server/handler/requesttype"
	"github.com/gofiber/fiber/v2"
)

func IsUsingApp(ctx *fiber.Ctx) bool {
	return GetAppVersion(ctx) != 0
}

func GetAppVersion(ctx *fiber.Ctx) int64 {
	ver := ctx.Locals("appVersion")
	if ver == nil {
		var req requesttype.AppHeader
		if err := ctx.ReqHeaderParser(&req); err != nil {
			return 0
		}
		return ctx.Locals("appVersion", req.Version).(int64)
	}
	return ver.(int64)
}
