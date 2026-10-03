package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

func RequestLogMiddleware(logger logrus.FieldLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestLogger := logger.WithFields(logrus.Fields{
			"method":    ctx.Method(),
			"path":      ctx.Path(),
			"ip":        ctx.IP(),
			"userAgent": string(ctx.Request().Header.Peek("user-agent")),
		})
		startTime := time.Now().UnixMilli()
		ctx.Locals("log", requestLogger)
		ret := ctx.Next()
		endTime := time.Now().UnixMilli()
		requestLogger.Infof("request success in %d ms", endTime-startTime)
		return ret
	}
}
