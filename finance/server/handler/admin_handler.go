package handler

import (
	"encoding/json"
	"fmt"
	"github.com/akmalfairuz/finance/server/handler/ctxhelper"
	"github.com/akmalfairuz/finance/server/handler/middleware"
	"github.com/akmalfairuz/finance/server/handler/requesttype"
	"github.com/akmalfairuz/finance/server/handler/responses"
	"github.com/akmalfairuz/finance/server/handler/responsetype"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AdminHandler struct {
	log                              logrus.FieldLogger
	authService                      *service.AuthService
	productService                   *service.ProductService
	purchaseService                  *service.PurchaseService
	notificationService              *service.NotificationService
	metaService                      *service.MetaService
	productExternalManagementService *service.ProductExternalManagementService
	supportService                   *service.SupportService
	userService                      *service.UserService
	kycService                       *service.UserKycService
	purchaseProcessorService         *service.PurchaseProcessorService
	analyticsService                 *service.AnalyticsService
}

func NewAdminHandler(log logrus.FieldLogger, authService *service.AuthService, productService *service.ProductService, notificationService *service.NotificationService, purchaseService *service.PurchaseService, metaService *service.MetaService, productExternalManagementService *service.ProductExternalManagementService, supportService *service.SupportService, kycService *service.UserKycService, userService *service.UserService, purchaseProcessorService *service.PurchaseProcessorService, analyticsService *service.AnalyticsService) *AdminHandler {
	return &AdminHandler{
		log:                              log,
		authService:                      authService,
		productService:                   productService,
		notificationService:              notificationService,
		purchaseService:                  purchaseService,
		metaService:                      metaService,
		productExternalManagementService: productExternalManagementService,
		supportService:                   supportService,
		kycService:                       kycService,
		userService:                      userService,
		purchaseProcessorService:         purchaseProcessorService,
	}
}

func (h *AdminHandler) Route(app *fiber.App) {
	authMiddleware := middleware.AuthMiddleware(h.authService)
	adminMiddleware := middleware.AdminMiddleware()

	group := app.Group("/admin", authMiddleware, adminMiddleware)

	group.Post("/announceClosing", h.handleAnnounceClosing)

	purchaseGroup := group.Group("/purchases")
	purchaseGroup.Get("/", h.handlePurchases)
	purchaseGroup.Post("/check", h.handleCheckPurchase)

	productGroup := group.Group("/products")
	productGroup.Get("/", h.handleProducts)
	productGroup.Get("/providers", h.handleProductProviders)
	productGroup.Post("/create", h.handleCreateProduct)
	productGroup.Post("/delete", h.handleDeleteProduct)
	productGroup.Post("/update", h.handleUpdateProduct)
	productGroup.Post("/importDigiflazz", h.handleImportProductFromDigiflazz)

	categoryGroup := productGroup.Group("/categories")
	categoryGroup.Post("/create", h.handleCreateCategory)
	categoryGroup.Post("/update", h.handleUpdateProductCategory)
	categoryGroup.Post("/delete", h.handleDeleteProductCategory)

	destinationGroup := productGroup.Group("/destinations")
	destinationGroup.Get("/", h.handleProductDestinations)
	destinationGroup.Post("/create", h.handleCreateProductDestination)
	destinationGroup.Post("/update", h.handleUpdateProductDestination)
	destinationGroup.Post("/delete", h.handleDeleteProductDestination)

	notificationGroup := group.Group("/notification")
	notificationGroup.Post("/push", h.handlePushNotification)

	supportGroup := group.Group("/support")
	supportGroup.Get("/tickets", h.handleSupportTickets)

	kycGroup := group.Group("/kyc")
	kycGroup.Get("/needToReview", h.handleKycNeedToReview)
	kycGroup.Post("/review", h.handleKycReview)
	kycGroup.Get("/detail", h.handleKycDetail)

	dashGroup := group.Group("/dashboard")
	dashGroup.Get("/summary", h.handleSummary)

	metaGroup := group.Group("/meta")
	metaGroup.Get("/", h.handleMeta)
	metaGroup.Post("/replace", h.handleReplaceMeta)

	//kycGroup := group.Group("/kyc")
	//kycGroup.Get("/", h.handle)

	//appGroup := group.Group("/app")
	//appGroup.Post("/updateVersion")
}

func (h *AdminHandler) handleAnnounceClosing(ctx *fiber.Ctx) error {
	go func() {
		if err := h.userService.AnnounceClosingToUser(); err != nil {
			h.log.Errorf("error sending closing email to user: %v", err)
		}
	}()
	return ctx.JSON(responses.M("Request sent"))
}

func (h *AdminHandler) handleCheckPurchase(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}

	purchase, err := h.purchaseService.GetPurchase(id)
	if err != nil {
		return err
	}

	if err := h.purchaseProcessorService.CheckStatus(purchase); err != nil {
		return err
	}

	return ctx.JSON(responses.M("Check request sent"))
}

func (h *AdminHandler) handleReplaceMeta(ctx *fiber.Ctx) error {
	var req requesttype.AdminReplaceMetaRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if err := h.metaService.Set(req.Name, req.Value); err != nil {
		return err
	}

	return ctx.JSON(responses.M("Replaced"))
}

func (h *AdminHandler) handleMeta(ctx *fiber.Ctx) error {
	metas, err := h.metaService.GetAll()
	if err != nil {
		return nil
	}

	resp := make([]responsetype.AdminMetaResponse, 0, len(metas))
	for _, meta := range metas {
		resp = append(resp, responsetype.AdminMetaResponse{
			Name:  meta.Name,
			Value: meta.Value,
		})
	}

	return ctx.JSON(resp)
}

func (h *AdminHandler) handleSummary(ctx *fiber.Ctx) error {
	// TODO
	return nil
}

func (h *AdminHandler) handleKycDetail(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}

	kyc, err := h.kycService.Get(id)
	if err != nil {
		return err
	}

	attachmentPrefix := h.kycService.GetAttachmentPrefix()

	return ctx.JSON(responsetype.NewAdminUserKycResponseFromModel(kyc, attachmentPrefix))
}

func (h *AdminHandler) handleKycNeedToReview(ctx *fiber.Ctx) error {
	kycs, err := h.kycService.GetAllPending()
	if err != nil {
		return nil
	}

	attachmentPrefix := h.kycService.GetAttachmentPrefix()

	resp := make([]responsetype.AdminUserKycResponse, 0, len(kycs))
	for _, kyc := range kycs {
		resp = append(resp, responsetype.NewAdminUserKycResponseFromModel(kyc, attachmentPrefix))
	}

	return ctx.JSON(resp)
}

func (h *AdminHandler) handleKycReview(ctx *fiber.Ctx) error {
	var req requesttype.AdminKycReviewRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	kyc, err := h.kycService.Get(req.ID)
	if err != nil {
		return err
	}

	if kyc.Status != 0 {
		return err
	}

	status := model.UserKycStatusFailed
	if req.Accepted {
		status = model.UserKycStatusSuccess
	}
	if err := h.kycService.UpdateStatus(kyc, status); err != nil {
		return err
	}

	return ctx.JSON(responses.M("Reviewed"))
}

func (h *AdminHandler) handleImportProductFromDigiflazz(ctx *fiber.Ctx) error {
	var req requesttype.AdminImportProductFromDigiflazzRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	importedProducts, err := h.productExternalManagementService.ImportFromDigiflazz(req.ToCategoryID, req.ProductKind, req.DestinationID, req.AddPrice, req.PrefixProductName, req.Brand, req.Type, req.Category)
	if err != nil {
		return err
	}

	return ctx.JSON(responses.M(fmt.Sprintf("Imported %d products", importedProducts)))
}

func (h *AdminHandler) handleSupportTickets(ctx *fiber.Ctx) error {
	tickets, err := h.supportService.GetSupportTickets()
	if err != nil {
		return err
	}
	resp := make([]responsetype.SupportTicketResponse, 0, len(tickets))
	for _, ticket := range tickets {
		resp = append(resp, responsetype.NewSupportTicketResponseFromModel(ticket))
	}
	return ctx.JSON(resp)
}

func (h *AdminHandler) handlePurchases(ctx *fiber.Ctx) error {
	var req requesttype.AdminPurchasesRequest
	if err := ctxhelper.BindQuery(ctx, &req); err != nil {
		return err
	}

	purchases, err := h.purchaseService.SearchPurchases(model.SearchPurchases{
		FromDate:         req.StartDate,
		ToDate:           req.EndDate,
		DestinationQuery: req.DestinationQuery,
	})
	if err != nil {
		return err
	}
	resp := make([]responsetype.AdminPurchaseResponse, 0)
	for _, purchase := range purchases {
		res, err := responsetype.NewAdminPurchaseResponseFromModel(purchase)
		if err != nil {
			return err
		}
		resp = append(resp, res)
	}
	return ctx.JSON(resp)
}

func (h *AdminHandler) handleProductDestinations(ctx *fiber.Ctx) error {
	destinations, err := h.productService.GetProductDestinations()
	if err != nil {
		return err
	}
	resp := make([]responsetype.ProductDestinationResponse, 0)
	for _, destination := range destinations {
		resp = append(resp, responsetype.NewProductDestinationResponseFromModel(destination, true))
	}
	return ctx.JSON(resp)
}

func (h *AdminHandler) handleCreateProductDestination(ctx *fiber.Ctx) error {
	var req requesttype.AdminCreateOrUpdateDestinationRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	if err := h.productService.CreateProductDestination(req.Model()); err != nil {
		return err
	}
	return ctx.JSON(&responsetype.MessageResponse{Message: "Product destination created"})
}

func (h *AdminHandler) handleUpdateProductDestination(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}
	var req requesttype.AdminCreateOrUpdateDestinationRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}
	if err := h.productService.UpdateProductDestination(id, req.Model()); err != nil {
		return err
	}
	return ctx.JSON(&responsetype.MessageResponse{Message: "Product destination updated"})
}

func (h *AdminHandler) handleDeleteProductDestination(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}
	if err := h.productService.DeleteProductDestination(id); err != nil {
		return err
	}
	return ctx.JSON(&responsetype.MessageResponse{Message: "Product destination deleted"})
}

func (h *AdminHandler) handleUpdateProduct(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}

	var req requesttype.AdminCreateOrUpdateProductRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if err := h.productService.UpdateProduct(id, req.Model()); err != nil {
		return err
	}

	return ctx.JSON(&responsetype.MessageResponse{Message: "Product updated"})
}

func (h *AdminHandler) handleDeleteProduct(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}

	if err := h.productService.DeleteProduct(id); err != nil {
		return err
	}

	return ctx.JSON(&responsetype.MessageResponse{Message: "Product deleted"})
}

func (h *AdminHandler) handleDeleteProductCategory(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}

	if err := h.productService.DeleteProductCategory(id); err != nil {
		return err
	}

	return ctx.JSON(&responsetype.MessageResponse{Message: "Product Category deleted"})
}

func (h *AdminHandler) handleProducts(ctx *fiber.Ctx) error {
	resp := make([]responsetype.AdminProductResponse, 0)

	products, err := h.productService.GetAllProducts()
	if err != nil {
		return err
	}

	for _, product := range products {
		resp = append(resp, responsetype.AdminProductResponse{
			ID:           product.ID,
			CategoryId:   product.CategoryID,
			CategoryName: product.CategoryName,
			Sku:          product.Sku,
			SkuPriority:  product.SkuPriority,
			ExternalSku:  product.ExternalSku,
			Provider:     product.Provider,
			Type:         product.Type,
			IsAvailable:  product.IsAvailable,
			Status:       product.StatusCode(),
			ImageUrl:     product.ImageUrl,
			Name:         product.Name,
			Kind:         product.Kind,
			Description:  product.Description,
			Destination: responsetype.ProductDestinationResponse{
				ID:          product.DestinationType,
				Name:        product.DestinationName,
				Description: product.DestinationDescription,
			},
			PurchaseNote:        product.PurchaseNote,
			Price:               product.Price,
			BillAdmin:           product.BillAdmin,
			MaxWholesalePrice:   product.MaxWholesalePrice,
			WholesalePrice:      product.WholesalePrice,
			BeforeDiscountPrice: product.BeforeDiscountPrice,
			CutOffStart:         product.CutOffStart,
			CutOffEnd:           product.CutOffEnd,
			ProofParser:         product.ProofParser,
		})
	}

	return ctx.JSON(&resp)
}

func (h *AdminHandler) handleProductProviders(ctx *fiber.Ctx) error {
	return ctx.JSON(model.ProductProviders)
}

func (h *AdminHandler) handleCreateCategory(ctx *fiber.Ctx) error {
	var req requesttype.AdminCreateOrUpdateCategoryRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if err := h.productService.CreateCategory(req.Model()); err != nil {
		return err
	}

	return nil
}

func (h *AdminHandler) handleUpdateProductCategory(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}

	var req requesttype.AdminCreateOrUpdateCategoryRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if err := h.productService.UpdateProductCategory(id, req.Model()); err != nil {
		return err
	}

	return ctx.JSON(&responsetype.MessageResponse{Message: "Product Category updated"})
}

// handleCreateProduct
func (h *AdminHandler) handleCreateProduct(ctx *fiber.Ctx) error {
	var req requesttype.AdminCreateOrUpdateProductRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	if err := h.productService.CreateProduct(req.Model()); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(&responsetype.MessageResponse{Message: "Product Created"})
}

func (h *AdminHandler) handlePushNotification(ctx *fiber.Ctx) error {
	var req requesttype.AdminPushNotificationRequest
	if err := ctxhelper.BindBody(ctx, &req); err != nil {
		return err
	}

	var data map[string]string
	if req.Data != "" {
		if err := json.Unmarshal([]byte(req.Data), &data); err != nil {
			return err
		}
	}

	if err := h.notificationService.PushNotificationToUser(req.UserID, &provider.PushNotificationOptions{
		Notification: &provider.NotificationOptions{
			Title:    req.Title,
			Body:     req.Body,
			ImageURL: req.ImageURL,
		},
		Data: data,
	}); err != nil {
		return err
	}

	return ctx.JSON(&responsetype.MessageResponse{Message: "Notification sent"})
}
