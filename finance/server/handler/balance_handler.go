package handler

import (
	"github.com/akmalfairuz/finance/server/handler/ctxhelper"
	"github.com/akmalfairuz/finance/server/handler/middleware"
	"github.com/akmalfairuz/finance/server/handler/requesttype"
	"github.com/akmalfairuz/finance/server/handler/responses"
	"github.com/akmalfairuz/finance/server/handler/responsetype"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type BalanceHandler struct {
	authService        *service.AuthService
	paymentService     *service.PaymentService
	transactionService *service.TransactionService
	userService        *service.UserService
	transferService    *service.TransferService
}

func NewBalanceHandler(authService *service.AuthService, paymentService *service.PaymentService, transactionService *service.TransactionService, userService *service.UserService, transferService *service.TransferService) *BalanceHandler {
	return &BalanceHandler{
		authService:        authService,
		paymentService:     paymentService,
		transactionService: transactionService,
		userService:        userService,
		transferService:    transferService,
	}
}

func (h *BalanceHandler) Route(app *fiber.App) {
	group := app.Group("/balance", middleware.AuthMiddleware(h.authService))
	group.Get("/", h.handleBalance)
	group.Get("/transactions", h.handleTransactions)

	transferGroup := group.Group("/transfer")
	transferGroup.Post("/", h.handleTransfer)
	transferGroup.Get("/checkDestination", h.handleTransferCheckDestination)
	transferGroup.Get("/history", h.handleTransferHistory)
	transferGroup.Get("/destinations", h.handleTransferDestinations)

	depositGroup := group.Group("/deposit")
	depositGroup.Post("/create", h.handleDepositCreate)
	depositGroup.Post("/cancel", h.handleDepositCancel)
	depositGroup.Get("/history", h.handleDepositHistory)
	depositGroup.Get("/detail", h.handleDepositDetail)
	depositGroup.Get("/paymentMethods", h.handlePaymentMethods)
}

func (h *BalanceHandler) handleTransferCheckDestination(ctx *fiber.Ctx) error {
	var req requesttype.BalanceTransferToUserCheckDestinationRequest
	if err := ctxhelper.BindQuery(ctx, &req); err != nil {
		return err
	}

	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}

	return ctx.JSON(&responsetype.BalanceTransferUserInfoResponse{
		Username: user.Name,
		Name:     user.DisplayName,
		Email:    user.Email,
	})
}

func (h *BalanceHandler) handleTransfer(ctx *fiber.Ctx) error {
	var req requesttype.BalanceTransferRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	receiverUser, err := h.userService.GetUserByEmail(req.DestinationEmail)
	if err != nil {
		return err
	}

	senderUser, err := h.userService.GetUser(ctxhelper.GetUserID(ctx))
	if err != nil {
		return err
	}

	if receiverUser.ID == senderUser.ID {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.M("You can't transfer to yourself"))
	}

	if !ctxhelper.ValidatePin(senderUser.PinHash, req.Pin) {
		return errortype.ErrInvalidPin
	}

	transferId, err := h.transferService.TransferToUser(senderUser, receiverUser, req.Amount, req.Note)
	if err != nil {
		return err
	}

	transfer, err := h.transferService.GetTransferToUserByID(transferId)
	if err != nil {
		return err
	}

	resp := responsetype.NewBalanceTransferToUserHistoryResponseFromModel(senderUser.ID, transfer)

	return ctx.JSON(&resp)
}

func (h *BalanceHandler) handleTransferHistory(ctx *fiber.Ctx) error {
	resp := make([]responsetype.BalanceTransferToUserHistoryResponse, 0)

	userId := ctxhelper.GetUserID(ctx)
	transfers, err := h.transferService.GetTransferToUserByUserId(userId)
	if err != nil {
		return err
	}

	for _, transfer := range transfers {
		resp = append(resp, responsetype.NewBalanceTransferToUserHistoryResponseFromModel(userId, transfer))
	}

	return ctx.JSON(resp)
}

func (h *BalanceHandler) handleTransferDestinations(ctx *fiber.Ctx) error {
	destinations, err := h.transferService.GetTransferToUserDestinations(ctxhelper.GetUserID(ctx))
	if err != nil {
		return err
	}

	resp := make([]responsetype.BalanceTransferUserInfoResponse, 0)
	for _, destination := range destinations {
		resp = append(resp, responsetype.BalanceTransferUserInfoResponse{
			Username: destination.Name,
			Name:     destination.DisplayName,
			Email:    destination.Email,
		})
	}

	return ctx.JSON(resp)
}

func (h *BalanceHandler) handlePaymentMethods(ctx *fiber.Ctx) error {
	resp := make([]responsetype.PaymentMethodResponse, 0)

	for _, paymentMethod := range h.paymentService.GetPaymentMethods() {
		resp = append(resp, responsetype.PaymentMethodResponse{
			ID:          paymentMethod.ID(),
			Name:        paymentMethod.Name(),
			Description: paymentMethod.Description(),
			NoFee:       paymentMethod.NoFee(),
			MinAmount:   paymentMethod.MinAmount(),
			MaxAmount:   paymentMethod.MaxAmount(),
			ImageUrl:    paymentMethod.ImageUrl(),
		})
	}

	return ctx.JSON(&resp)
}

// handleDepositCancel
func (h *BalanceHandler) handleDepositCancel(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}

	userId := ctxhelper.GetUserID(ctx)
	deposit, err := h.paymentService.GetDeposit(id)
	if err != nil {
		return err
	}

	if deposit.UserID != userId {
		return fiber.ErrForbidden
	}

	if err := h.paymentService.CancelDeposit(deposit.ID); err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON("success")
}

func (h *BalanceHandler) handleTransactions(ctx *fiber.Ctx) error {
	var req requesttype.GetTransactionsRequest
	if err := ctxhelper.BindQuery(ctx, &req); err != nil {
		return err
	}

	userId := ctxhelper.GetUserID(ctx)
	transactions, err := h.transactionService.GetUserTransactions(userId, req.From, req.To, req.Limit)
	if err != nil {
		return err
	}

	transactionInfos := make([]responsetype.TransactionInfo, 0)
	for _, transaction := range transactions {
		transactionInfos = append(transactionInfos, responsetype.TransactionInfo{
			Id:            transaction.ID,
			Amount:        transaction.Amount,
			Description:   transaction.Description,
			BeforeBalance: transaction.BeforeBalance,
			AfterBalance:  transaction.AfterBalance,
			CreatedAt:     transaction.CreatedAt,
		})
	}

	return ctx.Status(http.StatusOK).JSON(transactionInfos)
}

func (h *BalanceHandler) handleDepositHistory(ctx *fiber.Ctx) error {
	var req requesttype.DepositHistoryRequest
	if err := ctxhelper.BindQuery(ctx, &req); err != nil {
		return err
	}

	userId := ctxhelper.GetUserID(ctx)

	deposits, err := h.paymentService.GetDepositByUserId(userId, req.FilterStatus)
	if err != nil {
		return err
	}

	depositInfos := make([]responsetype.DepositInfo, 0)
	for _, deposit := range deposits {
		depositInfos = append(depositInfos, responsetype.DepositInfo{
			Id:            deposit.ID,
			Amount:        deposit.Amount,
			Fee:           deposit.Fee,
			PaymentMethod: h.paymentService.GetPaymentMethodName(deposit.PaymentMethod),
			Status:        deposit.Status,
			CreatedAt:     deposit.CreatedAt,
		})
	}

	return ctx.Status(http.StatusOK).JSON(depositInfos)
}

func (h *BalanceHandler) handleDepositDetail(ctx *fiber.Ctx) error {
	userID := ctxhelper.GetUserID(ctx)

	var req requesttype.DepositDetailedRequest
	if err := ctxhelper.BindQuery(ctx, &req); err != nil {
		return err
	}

	deposit, err := h.paymentService.GetDeposit(req.ID)
	if err != nil {
		return err
	}

	if deposit.UserID != userID {
		return fiber.ErrForbidden
	}

	publicInfo, err := h.paymentService.GetDepositPublicInfo(deposit.ID)
	if err != nil {
		return err
	}

	tracks, err := h.paymentService.GetDepositTracks(deposit.ID)
	if err != nil {
		return err
	}

	tracks2 := make([]responsetype.DepositTrack, 0)
	for _, track := range tracks {
		tracks2 = append(tracks2, responsetype.DepositTrack{
			Description: track.Description,
			NewStatus:   track.NewStatus,
			CreatedAt:   track.CreatedAt,
		})
	}

	// TODO: refactor this
	paymentUrl := ""
	if publicInfoMap, ok := publicInfo.(map[string]any); ok {
		if paymentUrlAny, ok := publicInfoMap["paymentUrl"]; ok {
			if paymentUrlString, ok := paymentUrlAny.(string); ok {
				paymentUrl = paymentUrlString
			}
		}
	}

	return ctx.Status(http.StatusOK).JSON(&responsetype.DepositDetailedInfo{
		Id:            deposit.ID,
		Amount:        deposit.Amount,
		Fee:           deposit.Fee,
		PaymentMethod: h.paymentService.GetPaymentMethodName(deposit.PaymentMethod),
		PaymentData:   publicInfo,
		PaymentUrl:    paymentUrl,
		Description:   deposit.Description,
		Status:        deposit.Status,
		CreatedAt:     deposit.CreatedAt,
		ExpiredAt:     deposit.ExpiredAt,
		Tracks:        tracks2,
	})
}

func (h *BalanceHandler) handleDepositCreate(ctx *fiber.Ctx) error {
	userId := ctxhelper.GetUserID(ctx)

	var req requesttype.DepositCreateRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	user, err := h.userService.GetUser(userId)
	if err != nil {
		return err
	}

	deposit, err := h.paymentService.CreateDeposit(user, req.Amount, req.PaymentMethod)
	if err != nil {
		return err
	}

	publicInfo, err := h.paymentService.GetDepositPublicInfo(deposit.ID)
	if err != nil {
		return err
	}

	tracks, err := h.paymentService.GetDepositTracks(deposit.ID)
	if err != nil {
		return err
	}

	tracks2 := make([]responsetype.DepositTrack, 0)
	for _, track := range tracks {
		tracks2 = append(tracks2, responsetype.DepositTrack{
			Description: track.Description,
			NewStatus:   track.NewStatus,
			CreatedAt:   track.CreatedAt,
		})
	}

	// TODO: refactor this
	paymentUrl := ""
	if publicInfoMap, ok := publicInfo.(map[string]any); ok {
		if paymentUrlAny, ok := publicInfoMap["paymentUrl"]; ok {
			if paymentUrlString, ok := paymentUrlAny.(string); ok {
				paymentUrl = paymentUrlString
			}
		}
	}

	return ctx.Status(http.StatusOK).JSON(&responsetype.DepositDetailedInfo{
		Id:            deposit.ID,
		Amount:        deposit.Amount,
		Fee:           deposit.Fee,
		PaymentMethod: h.paymentService.GetPaymentMethodName(deposit.PaymentMethod),
		PaymentData:   publicInfo,
		PaymentUrl:    paymentUrl,
		Description:   deposit.Description,
		Status:        deposit.Status,
		CreatedAt:     deposit.CreatedAt,
		ExpiredAt:     deposit.ExpiredAt,
		Tracks:        tracks2,
	})
}

func (h *BalanceHandler) handleBalance(ctx *fiber.Ctx) error {
	u, err := h.userService.GetUser(ctxhelper.GetUserID(ctx))
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(u.Balance)
}
