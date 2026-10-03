package handler

import (
	"github.com/akmalfairuz/finance/server/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type TestingHandler struct {
	otpService *service.OTPService
}

func NewTestingHandler(otpService *service.OTPService) *TestingHandler {
	return &TestingHandler{
		otpService: otpService,
	}
}

func (h *TestingHandler) Route(r fiber.Router) {
	group := r.Group("/testing")

	group.Get("/otp", h.handleOTP)
}

func (h *TestingHandler) handleOTP(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		return err
	}

	otp, err := h.otpService.GetOTP(int64(id))
	if err != nil {
		return err
	}

	return ctx.JSON(&otp)
}
