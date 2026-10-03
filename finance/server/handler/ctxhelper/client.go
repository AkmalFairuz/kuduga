package ctxhelper

import (
	"github.com/akmalfairuz/finance/module/helper"
	"github.com/gofiber/fiber/v2"
)

func SafeUniqueAddress(ctx *fiber.Ctx) string {
	ip := ctx.IP()
	ver := helper.IPVersion(ip)
	if ver == 4 {
		return ip
	}
	return helper.GetIPV6Prefix(ip)
}
