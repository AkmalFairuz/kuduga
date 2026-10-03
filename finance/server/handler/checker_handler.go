package handler

import (
	"github.com/akmalfairuz/finance/server/handler/ctxhelper"
	"github.com/akmalfairuz/finance/server/handler/middleware"
	"github.com/akmalfairuz/finance/server/handler/requesttype"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/gofiber/fiber/v2"
)

type CheckerHandler struct {
	authService    *service.AuthService
	checkerService *service.CheckerService
}

func NewCheckerHandler(checkerService *service.CheckerService, authService *service.AuthService) *CheckerHandler {
	return &CheckerHandler{checkerService: checkerService, authService: authService}
}

func (h *CheckerHandler) Route(r fiber.Router) {
	auth := middleware.AuthMiddleware(h.authService)

	group := r.Group("/checker", auth)
	group.Get("/", h.handleChecker)
}

func (h *CheckerHandler) handleChecker(ctx *fiber.Ctx) error {
	var req requesttype.CheckerRequest
	if err := ctxhelper.BindQuery(ctx, &req); err != nil {
		return err
	}
	input := map[string]string{}
	for _, v := range req.Input {
		input[v.Key] = v.Value
	}
	res, err := h.checkerService.Check(req.ID, input)
	if err != nil {
		return err
	}
	return ctx.JSON(&res)
}
