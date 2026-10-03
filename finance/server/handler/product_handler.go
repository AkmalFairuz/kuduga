package handler

import (
	"github.com/akmalfairuz/finance/module/pointer"
	"github.com/akmalfairuz/finance/server/handler/ctxhelper"
	"github.com/akmalfairuz/finance/server/handler/requesttype"
	"github.com/akmalfairuz/finance/server/handler/responsetype"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type ProductHandler struct {
	authService    *service.AuthService
	productService *service.ProductService
	metaService    *service.MetaService
}

func NewProductHandler(productService *service.ProductService, metaService *service.MetaService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		metaService:    metaService,
	}
}

func (h *ProductHandler) Route(app *fiber.App) {
	group := app.Group("/product")

	group.Get("/", h.handleProducts)
	group.Get("/single", h.handleSingleProduct)
	group.Get("/category", h.handleCategory)
	group.Get("/categories", h.handleCategories)
	group.Get("/destinations", h.handleDestinations)
	group.Get("/destination", h.handleDestination)

	group.Get("/pulsaCategories", h.handlePulsaCategories)
}

func (h *ProductHandler) handleCategories(ctx *fiber.Ctx) error {
	var req requesttype.GetProductCategoriesRequest
	if err := ctxhelper.BindQuery(ctx, &req); err != nil {
		return err
	}

	var categories []model.ProductCategory
	var err error

	if req.ParentID == 0 {
		categories, err = h.productService.GetAllProductCategoriesWithChildren(req.IsNested)
	} else {
		categories, err = h.productService.GetProductCategoryChildren(req.ParentID, req.IsNested)
	}

	if err != nil {
		return err
	}

	res := make([]responsetype.ProductCategoryResponse, 0, len(categories))
	for _, category := range categories {
		res = append(res, h.createProductCategoryResponse(category))
	}
	return ctx.Status(fiber.StatusOK).JSON(&res)
}

func (h *ProductHandler) handleCategory(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}

	productCategory, err := h.productService.GetProductCategory(id)
	if err != nil {
		return err
	}

	return ctx.JSON(h.createProductCategoryResponse(productCategory))
}

func (h *ProductHandler) createProductCategoryResponse(category model.ProductCategory) responsetype.ProductCategoryResponse {
	children := make([]responsetype.ProductCategoryResponse, 0)
	if len(category.Children) > 0 {
		children = make([]responsetype.ProductCategoryResponse, 0, len(category.Children))
		for _, child := range category.Children {
			children = append(children, h.createProductCategoryResponse(child))
		}
	}
	return responsetype.ProductCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		ParentID:    category.ParentID,
		IconUrl:     category.IconUrl,
		ProductType: category.ProductType,
		Description: category.Description,
		Children:    children,
		Meta:        category.Meta(),
	}
}

func (h *ProductHandler) handleProducts(ctx *fiber.Ctx) error {
	var products []model.Product
	var err error

	categoryId := ctx.Query("categoryId")
	if categoryId == "" {
		products, err = h.productService.GetAllProducts()
	} else {
		var categoryId2 int
		categoryId2, err = strconv.Atoi(categoryId)
		if err != nil {
			return err
		}
		products, err = h.productService.GetProductByCategory(int64(categoryId2))
	}

	if err != nil {
		return err
	}

	filteredProducts := filterProducts(products)

	res := make([]responsetype.ProductResponse, 0, len(filteredProducts))
	for _, product := range filteredProducts {
		res = append(res, responsetype.NewProductResponseFromModel(product))
	}
	return ctx.Status(fiber.StatusOK).JSON(&res)
}

func filterProducts(products []model.Product) []model.Product {
	marker := map[string]model.Product{}

	for _, product := range products {
		product2, ok := marker[product.Sku]
		if !ok {
			marker[product.Sku] = product
			continue
		}
		isProductAvailable := product.IsAvailable && !product.IsMaintenance()
		isProduct2Available := product2.IsAvailable && !product2.IsMaintenance()
		isBothAvailable := isProductAvailable && isProduct2Available
		if product2.SkuPriority == product.SkuPriority {
			if (isBothAvailable && product2.Price > product.Price) || (!isProduct2Available && isProductAvailable) {
				marker[product.Sku] = product
			}
			continue
		}
		if (isBothAvailable && product2.SkuPriority < product.SkuPriority) || (!isProduct2Available && isProductAvailable) {
			marker[product.Sku] = product
		}
	}

	ret := make([]model.Product, 0, len(marker))
	for _, product := range marker {
		ret = append(ret, product)
	}

	return ret
}

func (h *ProductHandler) handlePulsaCategories(ctx *fiber.Ctx) error {
	pulsaPrefix, err := h.metaService.Get("pulsa_prefix")
	if err != nil {
		return err
	}

	resp := map[string]string{}

	for _, prefix := range strings.Split(pulsaPrefix, "\n") {
		prefix2 := strings.Split(prefix, "=")
		if len(prefix2) != 2 {
			continue
		}
		resp[prefix2[0]] = prefix2[1]
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (h *ProductHandler) handleDestinations(ctx *fiber.Ctx) error {
	resp := make([]responsetype.ProductDestinationResponse, 0)

	destinations, err := h.productService.GetProductDestinations()
	if err != nil {
		return err
	}

	for _, destination := range destinations {
		resp = append(resp, responsetype.NewProductDestinationResponseFromModel(destination, false))
	}

	return ctx.JSON(&resp)
}

func (h *ProductHandler) handleDestination(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}
	destination, err := h.productService.GetProductDestination(id)
	if err != nil {
		return err
	}
	resp := responsetype.NewProductDestinationResponseFromModel(destination, false)
	return ctx.JSON(&resp)
}

func (h *ProductHandler) handleSingleProduct(ctx *fiber.Ctx) error {
	id, err := ctxhelper.ParseID(ctx)
	if err != nil {
		return err
	}
	product, err := h.productService.GetProduct(id)
	return ctx.JSON(pointer.Make(responsetype.NewProductResponseFromModel(product)))
}
