package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/akmalfairuz/finance/server/handler/ctxhelper"
	"github.com/akmalfairuz/finance/server/handler/middleware"
	"github.com/akmalfairuz/finance/server/handler/requesttype"
	"github.com/akmalfairuz/finance/server/handler/responses"
	"github.com/akmalfairuz/finance/server/handler/responsetype"
	"github.com/akmalfairuz/finance/server/handler/wshelper"
	"github.com/akmalfairuz/finance/server/handler/xhelper"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type SupportHandler struct {
	authService    *service.AuthService
	userService    *service.UserService
	supportService *service.SupportService
}

func NewSupportHandler(supportService *service.SupportService, authService *service.AuthService, userService *service.UserService) *SupportHandler {
	return &SupportHandler{supportService: supportService, authService: authService, userService: userService}
}

func (h *SupportHandler) Route(r fiber.Router) {
	auth := middleware.AuthMiddleware(h.authService)

	group := r.Group("/support")

	ticketGroup := group.Group("/ticket", auth)
	ticketGroup.Get("/", h.handleTickets)
	ticketGroup.Get("/detail", h.handleTicketDetail)
	ticketGroup.Post("/message", h.handleCreateTicketMessage)
	ticketGroup.Post("/create", h.handleCreateTicket)
	ticketGroup.Post("/close", h.handleCloseTicket)
	ticketGroup.Get("/categories", h.handleTicketCategories)

	// TODO: this is a bad code, the route path MUST not hardcoded like this
	// because browser can't pass authentication to websocket endpoint
	r.Get("/ws-support/listenTicketUpdate", h.handlePreListenSupportTicketUpdate, websocket.New(h.handleListenSupportTicketUpdate))
}

func (h *SupportHandler) handleCloseTicket(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}

	ticket, err := h.supportService.GetSupportTicketByID(id)
	if err != nil {
		return err
	}

	if !ctxhelper.IsAdmin(ctx) && ticket.UserID != ctxhelper.GetUserID(ctx) {
		return ctx.Status(fiber.StatusForbidden).JSON(responses.M("You don't have access to this ticket"))
	}

	if err := h.supportService.CloseTicket(id); err != nil {
		return err
	}

	return ctx.JSON(responses.M("Ticket closed"))
}

func (h *SupportHandler) handlePreListenSupportTicketUpdate(ctx *fiber.Ctx) error {
	ticketId, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}
	ctx.Locals("ticketId", ticketId)
	if websocket.IsWebSocketUpgrade(ctx) {
		return ctx.Next()
	}
	return ctx.Status(fiber.StatusUpgradeRequired).JSON(responses.M("Websocket required"))
}

func (h *SupportHandler) handleListenSupportTicketUpdate(conn *websocket.Conn) {
	initMsg, err := wshelper.InitialCheck(conn)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return
	}
	tok := xhelper.GetToken(h.authService, initMsg.Authorization)
	if tok == nil {
		_ = conn.WriteJSON(responses.M("Invalid token"))
		return
	}

	ticketId := conn.Locals("ticketId").(int64)

	ticket, err := h.supportService.GetSupportTicketByID(ticketId)
	if err != nil {
		return
	}

	if tok.UserID != ticket.UserID && tok.Role != model.RoleAdmin {
		_ = conn.WriteJSON(responses.M("Forbidden"))
		return
	}

	ctx, cancel := context.WithCancel(context.TODO())
	conn.SetCloseHandler(func(code int, text string) error {
		cancel()
		return nil
	})

	if err := h.supportService.ListenSupportTicketUpdate(ctx, ticketId, func() error {
		return conn.WriteJSON("1")
	}); err != nil {
		cancel()
	}
}

func (h *SupportHandler) handleTicketCategories(ctx *fiber.Ctx) error {
	return ctx.JSON(model.SupportTicketCategories)
}

func (h *SupportHandler) handleTicketDetail(ctx *fiber.Ctx) error {
	var req requesttype.GetTicketDetailRequest
	if err := ctxhelper.BindQuery(ctx, &req); err != nil {
		return err
	}

	ticket, err := h.supportService.GetSupportTicketByID(req.ID)
	if err != nil {
		return err
	}

	if !ctxhelper.IsAdmin(ctx) && ticket.UserID != ctxhelper.GetUserID(ctx) {
		return ctx.Status(fiber.StatusForbidden).JSON(responses.M("You don't have an access to this ticket"))
	}

	messages, err := h.supportService.GetSupportTicketMessagesAfterID(req.ID, req.MessageAfterID)
	if err != nil {
		return err
	}

	msgResp := make([]responsetype.SupportTicketMessageResponse, 0)
	for _, message := range messages {
		msgResp = append(msgResp, responsetype.NewSupportTicketMessageResponseFromModel(message, h.supportService.GetAttachmentPrefix()))
	}

	resp := responsetype.NewSupportTicketResponseFromModel(ticket)
	resp.Messages = msgResp

	return ctx.JSON(&resp)
}

func (h *SupportHandler) handleCreateTicketMessage(ctx *fiber.Ctx) error {
	var req requesttype.CreateSupportTicketMessageRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	ticket, err := h.supportService.GetSupportTicketByID(req.TicketID)
	if err != nil {
		return err
	}

	user, err := h.userService.GetUser(ctxhelper.GetUserID(ctx))
	if err != nil {
		return err
	}

	if !ctxhelper.IsAdmin(ctx) && ticket.UserID != user.ID {
		return ctx.Status(fiber.StatusForbidden).JSON(responses.M("You don't have an access to this ticket"))
	}

	files, err := h.readFiles(ctx)
	if err != nil {
		return err
	}

	// TODO: this is dumb, move this to admin_handler.go
	role := model.SupportTicketUserRole
	if ctxhelper.IsAdmin(ctx) && req.IsCustomerService {
		role = model.SupportTicketCustomerServiceRole
	}

	if _, err := h.supportService.CreateSupportTicketMessage(&model.CreateSupportTicketMessage{
		TicketID:    req.TicketID,
		Author:      user.Name,
		Role:        role,
		Message:     req.Message,
		Attachments: files,
	}); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.M("Message created"))
}

func (h *SupportHandler) handleCreateTicket(ctx *fiber.Ctx) error {
	var req requesttype.CreateSupportTicketRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	user, err := h.userService.GetUser(ctxhelper.GetUserID(ctx))
	if err != nil {
		return err
	}

	files, err := h.readFiles(ctx)
	if err != nil {
		return err
	}

	ticketId, err := h.supportService.CreateSupportTicket(model.CreateSupportTicket{
		UserID:   user.ID,
		Category: req.CategoryID,
	}, &model.CreateSupportTicketMessage{
		Author:      user.Name,
		Role:        model.SupportTicketUserRole,
		Message:     req.Message,
		Attachments: files,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(&responsetype.CreateSupportTicketResponse{
		TicketID: ticketId,
	})
}

func (h *SupportHandler) handleTickets(ctx *fiber.Ctx) error {
	userId := ctxhelper.GetUserID(ctx)
	tickets, err := h.supportService.GetSupportTicketsByUserID(userId)
	if err != nil {
		return err
	}

	resp := make([]responsetype.SupportTicketResponse, 0)
	for _, ticket := range tickets {
		resp = append(resp, responsetype.NewSupportTicketResponseFromModel(ticket))
	}

	return ctx.JSON(resp)
}

func (h *SupportHandler) readFiles(ctx *fiber.Ctx) ([]*model.File, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		if err == fasthttp.ErrNoMultipartForm {
			return []*model.File{}, nil
		}
		return nil, err
	}

	files := make([]*model.File, 0)
	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			if fileHeader.Size >= 1e7 {
				return nil, errors.New("file must under 10mb")
			}
			file, err := fileHeader.Open()
			if err != nil {
				return nil, err
			}
			files = append(files, &model.File{Name: fileHeader.Filename, Bytes: file})
		}
	}
	return files, nil
}
