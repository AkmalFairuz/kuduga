package handler

import (
	"encoding/json"
	"github.com/akmalfairuz/finance/server/handler/responsetype"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/gofiber/fiber/v2"
)

type AppHandler struct {
	metaService *service.MetaService
}

func NewAppHandler(metaService *service.MetaService) *AppHandler {
	return &AppHandler{metaService: metaService}
}

func (h *AppHandler) Route(r fiber.Router) {
	group := r.Group("/app")
	group.Get("/layoutInfo", h.handleLayoutInfo)
	group.Get("/bannerInfo", h.handleBannerInfo)
	group.Get("/privacyPolicy", h.handlePrivacyPolicy)
	group.Get("/termsAndConditions", h.handleTermsAndConditions)
	group.Get("/version", h.handleVersion)
}

func (h *AppHandler) handleVersion(ctx *fiber.Ctx) error {
	ret, err := h.metaService.Get("app_version")
	if err != nil {
		return err
	}
	var resp responsetype.AppVersionResponse
	if err := json.Unmarshal([]byte(ret), &resp); err != nil {
		return err
	}
	return ctx.JSON(&resp)
}

func (h *AppHandler) handleBannerInfo(ctx *fiber.Ctx) error {
	ret, err := h.metaService.Get("app_banner_info")
	if err != nil {
		return err
	}
	resp := make([]responsetype.AppBannerInfoResponse, 0)
	if err := json.Unmarshal([]byte(ret), &resp); err != nil {
		return err
	}
	return ctx.JSON(&resp)
}

func (h *AppHandler) handlePrivacyPolicy(ctx *fiber.Ctx) error {
	ret, err := h.metaService.Get("app_privacy_policy")
	if err != nil {
		return err
	}
	ctx.Response().Header.Set("Content-Type", "text/html")
	return ctx.SendString(ret)
}

func (h *AppHandler) handleTermsAndConditions(ctx *fiber.Ctx) error {
	ret, err := h.metaService.Get("app_terms_and_conditions")
	if err != nil {
		return err
	}
	ctx.Response().Header.Set("Content-Type", "text/html")
	return ctx.SendString(ret)
}

func (h *AppHandler) handleLayoutInfo(ctx *fiber.Ctx) error {
	ret, err := h.metaService.Get("app_layout_info")
	if err != nil {
		return err
	}
	var resp responsetype.AppLayoutInfoResponse
	if err := json.Unmarshal([]byte(ret), &resp); err != nil {
		return err
	}
	return ctx.JSON(&resp)
}
