package handler

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/module/helper"
	"github.com/akmalfairuz/finance/module/pointer"
	"github.com/akmalfairuz/finance/server/handler/ctxhelper"
	"github.com/akmalfairuz/finance/server/handler/middleware"
	"github.com/akmalfairuz/finance/server/handler/requesttype"
	"github.com/akmalfairuz/finance/server/handler/responses"
	"github.com/akmalfairuz/finance/server/handler/responsetype"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/gofiber/fiber/v2"
)

type PurchaseHandler struct {
	purchaseService *service.PurchaseService
	authService     *service.AuthService
	userService     *service.UserService
}

func NewPurchaseHandler(purchaseService *service.PurchaseService, authService *service.AuthService, userService *service.UserService) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseService: purchaseService,
		authService:     authService,
		userService:     userService,
	}
}

func (h *PurchaseHandler) Route(app *fiber.App) {
	auth := middleware.AuthMiddleware(h.authService)

	group := app.Group("/purchase")
	group.Post("/create", auth, h.handleCreatePurchase)
	group.Post("/bill", auth, h.handleBill)
	group.Post("/updateUserSellPrice", auth, h.handleUpdateUserSellPrice)
	group.Get("/status", auth, h.handleStatus)
	group.Get("/history", auth, h.handleHistory)
	group.Get("/stats", auth, h.handleStats)
	group.Get("/search", auth, h.handleSearch)
	group.Get("/publicHistory", h.handlePublicHistory)
}

func (h *PurchaseHandler) handleSearch(ctx *fiber.Ctx) error {
	var req requesttype.SearchPurchaseRequest
	if err := ctxhelper.BindQuery(ctx, &req); err != nil {
		return err
	}
	userId := ctxhelper.GetUserID(ctx)

	purchases, err := h.purchaseService.SearchPurchasesByUserId(userId, req.Query)
	if err != nil {
		return err
	}

	resp := make([]responsetype.PurchaseResponse, 0)
	for _, purchase := range purchases {
		p, err := h.toPurchaseResponse(purchase)
		if err != nil {
			return err
		}
		resp = append(resp, p)
	}

	return ctx.JSON(resp)
}

func (h *PurchaseHandler) handleStats(ctx *fiber.Ctx) error {
	userId := ctxhelper.GetUserID(ctx)

	resp := responsetype.UserPurchaseStatsWrapperResponse{}

	{
		start, end := helper.ThisMonthUTC7Range()
		stats, err := h.purchaseService.GetUserPurchaseStats(userId, start, end)
		if err != nil {
			return err
		}
		resp.ThisMonth = responsetype.NewUserPurchaseStatsResponse(stats)
	}
	{
		start, end := helper.PreviousMonthUTC7Range()
		stats, err := h.purchaseService.GetUserPurchaseStats(userId, start, end)
		if err != nil {
			return err
		}
		resp.PreviousMonth = responsetype.NewUserPurchaseStatsResponse(stats)
	}
	{
		start, end := helper.YesterdayUTC7Range()
		stats, err := h.purchaseService.GetUserPurchaseStats(userId, start, end)
		if err != nil {
			return err
		}
		resp.Yesterday = responsetype.NewUserPurchaseStatsResponse(stats)
	}
	{
		start, end := helper.TodayUTC7Range()
		stats, err := h.purchaseService.GetUserPurchaseStats(userId, start, end)
		if err != nil {
			return err
		}
		resp.Today = responsetype.NewUserPurchaseStatsResponse(stats)
	}

	return ctx.JSON(&resp)
}

func (h *PurchaseHandler) handlePublicHistory(ctx *fiber.Ctx) error {
	purchases, err := h.purchaseService.GetLatestPurchases(40)
	if err != nil {
		return err
	}
	resp := make([]responsetype.PublicPurchaseResponse, 0)
	for _, purchase := range purchases {
		resp = append(resp, responsetype.PublicPurchaseResponse{
			PurchaseId:          purchase.ID,
			ProductSku:          purchase.ProductSku,
			ProductName:         purchase.ProductName,
			ProductKind:         purchase.ProductKind,
			ProductType:         purchase.ProductType,
			ProductCategoryName: purchase.ProductCategoryName,
			Price:               purchase.Price,
			Status:              purchase.Status,
			Destination:         blurPurchaseDestination(purchase.Destination),
			CreatedAt:           purchase.CreatedAt,
		})
	}

	return ctx.JSON(resp)
}

func (h *PurchaseHandler) handleUpdateUserSellPrice(ctx *fiber.Ctx) error {
	var req requesttype.UpdatePurchaseUserSellPriceRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	purchase, err := h.purchaseService.GetPurchase(req.ID)
	if err != nil {
		if err == database.ErrNoRows {
			return errortype.ErrPurchaseNotFound
		}
		return err
	}
	if purchase.UserID == nil || *purchase.UserID != ctxhelper.GetUserID(ctx) {
		return errortype.ErrPurchaseNoAccess
	}
	if err := h.purchaseService.UpdatePurchaseUserSellPrice(purchase, req.UserSellPrice); err != nil {
		return err
	}
	return ctx.JSON(responses.M("User sell price updated"))
}

func (h *PurchaseHandler) handleHistory(ctx *fiber.Ctx) error {
	var req requesttype.PurchaseHistoryRequest
	if err := ctxhelper.BindQuery(ctx, &req); err != nil {
		return err
	}

	userId := ctxhelper.GetUserID(ctx)

	purchases, err := h.purchaseService.GetPurchasesByUserId(userId, model.GetUserPurchasesOptions{
		Status:   req.Status,
		BeforeID: req.BeforeID,
		Limit:    req.Limit,
	})
	if err != nil {
		return err
	}

	resp := make([]responsetype.PurchaseResponse, 0)
	for _, purchase := range purchases {
		p, err := h.toPurchaseResponse(purchase)
		if err != nil {
			return err
		}
		resp = append(resp, p)
	}

	return ctx.JSON(resp)
}

func (h *PurchaseHandler) handleStatus(ctx *fiber.Ctx) error {
	userId := ctxhelper.GetUserID(ctx)

	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}

	purchase, err := h.purchaseService.GetPurchase(id)
	if err != nil {
		if err == database.ErrNoRows {
			return errortype.ErrPurchaseNotFound
		}
		return err
	}

	if purchase.UserID == nil || *purchase.UserID != userId {
		return errortype.ErrPurchaseNoAccess
	}

	resp, err := h.toPurchaseResponse(purchase)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(pointer.Make(resp))
}

func (h *PurchaseHandler) handleCreatePurchase(ctx *fiber.Ctx) error {
	userId := ctxhelper.GetUserID(ctx)

	user, err := h.userService.GetUser(userId)
	if err != nil {
		return err
	}

	var req requesttype.CreatePurchaseRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if !ctxhelper.ValidatePin(user.PinHash, req.Pin) {
		return errortype.ErrInvalidPin
	}

	result, err := h.purchaseService.CreatePurchase(&model.CreatePurchaseOptions{
		UserID:          &userId,
		UserKycVerified: user.KycStatus,
		Email:           user.Email,
		Destination:     req.Destination(),
		ProductID:       req.ProductID,
	})
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(&responsetype.CreatePurchaseResponse{
		PurchaseId: result.PurchaseID,
	})
}

func (h *PurchaseHandler) toPurchaseResponse(purchase model.Purchase) (responsetype.PurchaseResponse, error) {
	destination, err := purchase.ReadableDestination()
	if err != nil {
		return responsetype.PurchaseResponse{}, err
	}
	extraData, err := purchase.ExtraData()
	if err != nil {
		return responsetype.PurchaseResponse{}, err
	}
	proof := purchase.Proof
	if proof != nil && purchase.ParsedProof != "" {
		proof = &purchase.ParsedProof
	}
	resp := responsetype.PurchaseResponse{
		PurchaseId:          purchase.ID,
		ProductId:           purchase.ProductID,
		ProductSku:          purchase.ProductSku,
		ProductName:         purchase.ProductName,
		ProductKind:         purchase.ProductKind,
		ProductType:         purchase.ProductType,
		ProductCategoryId:   purchase.ProductCategoryId,
		ProductCategoryName: purchase.ProductCategoryName,
		ProductDescription:  purchase.ProductDescription,
		Status:              purchase.Status,
		StatusText:          purchase.StatusText(),
		Destination:         destination,
		ExtraData:           extraData,
		Note:                purchase.Note,
		Proof:               proof,
		Price:               purchase.Price,
		UserSellPrice:       purchase.UserSellPrice,
		Fee:                 purchase.Fee,
		PaymentMethod:       purchase.PaymentMethod,
		PaymentExpiredAt:    purchase.PaymentExpiredAt,
		PaymentData:         purchase.PaymentData,
		CreatedAt:           purchase.CreatedAt,
	}
	if purchase.ProductType == model.ProductTypeDigitalPostpaid {
		billData, err := purchase.BillData()
		if err != nil {
			return responsetype.PurchaseResponse{}, err
		}
		resp.TotalBill = purchase.TotalBill
		resp.TotalBillFee = purchase.TotalBillFee
		resp.BillAdmin = &purchase.BillAdmin
		resp.BillData = billData
	}
	return resp, nil
}

func (h *PurchaseHandler) handleBill(ctx *fiber.Ctx) error {
	var req requesttype.BillRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	userId := ctxhelper.GetUserID(ctx)
	user, err := h.userService.GetUser(userId)
	if err != nil {
		return err
	}

	if req.Pay {
		if !ctxhelper.ValidatePin(user.PinHash, req.Pin) {
			return errortype.ErrInvalidPin
		}
	} else {
		result, err := h.purchaseService.CheckBill(&userId, user.Email, req.Destination(), req.ProductID)
		if err != nil {
			return err
		}

		return ctx.JSON(&responsetype.CheckBillResponse{
			ID:           result.ID,
			BillAmount:   result.BillAmount,
			Admin:        result.Admin,
			RealAdmin:    result.RealAdmin,
			TotalPrice:   result.CustomerPrice,
			CustomerName: result.CustomerName,
			Data:         result.Data,
			BillData:     result.BillData,
		})
	}

	bill, err := h.purchaseService.GetBillPrePurchase(req.BillID)
	if err != nil {
		return err
	}
	if bill.UserID != nil && *bill.UserID != userId {
		return ctx.Status(fiber.StatusBadRequest).JSON(responses.M("Invalid Bill ID"))
	}

	result, err := h.purchaseService.CreatePurchase(&model.CreatePurchaseOptions{
		UserID:          &userId,
		UserKycVerified: user.KycStatus,
		Email:           user.Email,
		Destination:     req.Destination(),
		ProductID:       req.ProductID,
		Bill:            &bill,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(&responsetype.CreatePurchaseResponse{PurchaseId: result.PurchaseID})
}

func blurPurchaseDestination(destination string) string {
	if len(destination) >= 11 {
		return destination[:7] + "****"
	}
	if len(destination) >= 8 {
		return destination[:5] + "****"
	}
	return "****"
}
