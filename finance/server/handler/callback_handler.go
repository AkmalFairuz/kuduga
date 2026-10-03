package handler

import (
	"github.com/akmalfairuz/finance/module/duitku"
	"github.com/akmalfairuz/finance/server/handler/ctxhelper"
	"github.com/akmalfairuz/finance/server/handler/requesttype"
	"github.com/akmalfairuz/finance/server/handler/responses"
	"github.com/akmalfairuz/finance/server/handler/responsetype"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/akmalfairuz/finance/server/service/paymentmethod"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type CallbackHandler struct {
	purchaseProcessor *service.PurchaseProcessorService
	paymentService    *service.PaymentService
	duitkuClient      *duitku.Client

	bcaMutasiApiKey string
}

func NewCallbackHandler(purchaseProcessor *service.PurchaseProcessorService, paymentService *service.PaymentService) *CallbackHandler {
	return &CallbackHandler{
		purchaseProcessor: purchaseProcessor,
		paymentService:    paymentService,
	}
}

func (h *CallbackHandler) SetBcaMutasiApiKey(apiKey string) {
	h.bcaMutasiApiKey = apiKey
}

func (h *CallbackHandler) SetDuitkuClient(client *duitku.Client) {
	h.duitkuClient = client
}

func (h *CallbackHandler) Route(app *fiber.App) {
	group := app.Group("/callback")

	paymentGroup := group.Group("/payment")
	paymentGroup.Post("/bca", h.handleBcaMutasiPayment)
	paymentGroup.Post("/duitku", h.handleDuitkuPayment)
	// backward compability
	group.Post("/duitku", h.handleDuitkuPayment)

	group.Post("/digiflazz", h.handleDigiflazz)
}

func (h *CallbackHandler) handleDigiflazz(ctx *fiber.Ctx) error {
	if err := h.purchaseProcessor.HandleCallback(provider.DigiflazzPurchaseProviderName, ctx.Context()); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(&responsetype.MessageResponse{Message: "OK"})
}

func (h *CallbackHandler) handleBcaMutasiPayment(ctx *fiber.Ctx) error {
	var req requesttype.PaymentCallbackBcaMutasiRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if h.bcaMutasiApiKey == "" || req.ApiKey != h.bcaMutasiApiKey {
		return ctx.Status(http.StatusUnauthorized).JSON(responses.M("Invalid api key"))
	}

	if req.Status == 1 { // Success
		if err := h.paymentService.HandlePaymentSuccess(req.RefId, paymentmethod.BankBcaID); err != nil {
			return err
		}
	}

	ctx.Status(http.StatusOK)
	return nil
}

func (h *CallbackHandler) handleDuitkuPayment(ctx *fiber.Ctx) error {
	result, err := h.duitkuClient.HandleCallback(ctx)
	if err != nil {
		return err
	}

	switch result.ResultCode {
	case duitku.SuccessResultCode:
		if err := h.paymentService.HandlePaymentSuccess(result.MerchantOrderID, paymentmethod.FromDuitkuCode(result.PaymentCode)); err != nil {
			return err
		}
	case duitku.FailedResultCode:
		if err := h.paymentService.HandlePaymentFailed(result.MerchantOrderID, paymentmethod.FromDuitkuCode(result.PaymentCode)); err != nil {
			return err
		}
	}

	return ctx.JSON(responses.M("handled"))
}
