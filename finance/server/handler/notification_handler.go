package handler

import (
	"github.com/akmalfairuz/finance/server/handler/ctxhelper"
	"github.com/akmalfairuz/finance/server/handler/middleware"
	"github.com/akmalfairuz/finance/server/handler/requesttype"
	"github.com/akmalfairuz/finance/server/handler/responses"
	"github.com/akmalfairuz/finance/server/handler/responsetype"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	authService         *service.AuthService
	notificationService *service.NotificationService
}

func NewNotificationHandler(authService *service.AuthService, notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		authService:         authService,
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) Route(router fiber.Router) {
	auth := middleware.AuthMiddleware(h.authService)
	group := router.Group("/notification", auth)
	group.Post("/updateFcmToken", h.handleUpdateFcmToken)
	group.Get("/", h.handleGet)
	group.Post("/markRead", h.handleMarkRead)
}

func (h *NotificationHandler) handleMarkRead(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}
	userId := ctxhelper.GetUserID(ctx)
	if err := h.notificationService.MarkAsReadNotification(userId, id); err != nil {
		return err
	}
	return ctx.JSON(responses.M("Marked as read"))
}

func (h *NotificationHandler) handleGet(ctx *fiber.Ctx) error {
	userId := ctxhelper.GetUserID(ctx)
	notifications, err := h.notificationService.GetNotificationsByUserId(userId)
	if err != nil {
		return err
	}
	resp := make([]responsetype.NotificationResponse, 0)
	for _, notification := range notifications {
		data, err := notification.Data()
		if err != nil {
			return err
		}
		resp = append(resp, responsetype.NotificationResponse{
			ID:          notification.ID,
			Title:       notification.Title,
			Description: notification.Description,
			HasRead:     notification.HasRead,
			RefID:       notification.RefID,
			Type:        notification.Type,
			Data:        data,
			CreatedAt:   notification.CreatedAt,
		})
	}

	return ctx.JSON(resp)
}

func (h *NotificationHandler) handleUpdateFcmToken(ctx *fiber.Ctx) error {
	var req requesttype.UpdateFcmTokenRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	tok := ctxhelper.GetToken(ctx)

	if tok.DeviceFcmToken == nil || *tok.DeviceFcmToken != req.FcmToken {
		tokenID := tok.ID
		if err := h.authService.SetDeviceFcmToken(tokenID, req.FcmToken); err != nil {
			return err
		}
		if err := h.notificationService.InitializeFcmToken(req.FcmToken); err != nil {
			return err
		}
	}

	return ctx.JSON(responsetype.MessageResponse{Message: "FCM token updated"})
}
