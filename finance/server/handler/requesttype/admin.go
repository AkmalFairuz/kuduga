package requesttype

import "github.com/akmalfairuz/finance/server/model"

type AdminCreateOrUpdateProductRequest struct {
	CategoryId          int64   `form:"categoryId" validate:"required"`
	Sku                 string  `form:"sku" validate:"required"`
	SkuPriority         *int    `form:"skuPriority"`
	ExternalSku         string  `form:"externalSku" validate:"required"`
	Provider            string  `form:"provider" validate:"required"`
	Type                *int    `form:"type" validate:"required"`
	ImageUrl            *string `form:"imageUrl"`
	CutOffStart         *string `form:"cutOffStart"`
	CutOffEnd           *string `form:"cutOffEnd"`
	IsAvailable         *bool   `form:"isAvailable" validate:"required"`
	Name                string  `form:"name" validate:"required"`
	Kind                string  `form:"kind"`
	PurchaseNote        string  `form:"purchaseNote"`
	Description         string  `form:"description" validate:"max=5000"`
	DestinationType     int64   `form:"destinationType" validate:"required"`
	Price               int64   `form:"price" validate:"required"`
	BillAdmin           int64   `form:"billAdmin"`
	MaxWholesalePrice   int64   `form:"maxWholesalePrice" validate:"required"`
	WholesalePrice      int64   `form:"wholesalePrice" validate:"required"`
	BeforeDiscountPrice int64   `form:"beforeDiscountPrice"`
	ProofParser         string  `form:"proofParser"`
}

func (req AdminCreateOrUpdateProductRequest) Model() model.CreateOrUpdateProduct {
	return model.CreateOrUpdateProduct{
		CategoryID:          req.CategoryId,
		Sku:                 req.Sku,
		SkuPriority:         req.SkuPriority,
		ExternalSku:         req.ExternalSku,
		Provider:            req.Provider,
		Name:                req.Name,
		Kind:                &req.Kind,
		PurchaseNote:        &req.PurchaseNote,
		Description:         &req.Description,
		Type:                *req.Type,
		ImageUrl:            req.ImageUrl,
		CutOffStart:         req.CutOffStart,
		CutOffEnd:           req.CutOffEnd,
		DestinationType:     req.DestinationType,
		IsAvailable:         req.IsAvailable,
		Price:               req.Price,
		BillAdmin:           &req.BillAdmin,
		MaxWholesalePrice:   req.MaxWholesalePrice,
		WholesalePrice:      req.WholesalePrice,
		BeforeDiscountPrice: &req.BeforeDiscountPrice,
		ProofParser:         &req.ProofParser,
	}
}

type AdminCreateOrUpdateCategoryRequest struct {
	Name        string  `form:"name" validate:"required"`
	ParentId    *int64  `form:"parentId"`
	Description *string `form:"description"`
	IconUrl     *string `form:"iconUrl"`
	ProductType int     `form:"productType" validate:"required"`
	Meta        *string `form:"meta"`
}

func (req AdminCreateOrUpdateCategoryRequest) Model() model.CreateOrUpdateProductCategory {
	return model.CreateOrUpdateProductCategory{
		Name:        req.Name,
		ParentID:    &req.ParentId,
		Description: req.Description,
		IconUrl:     req.IconUrl,
		ProductType: req.ProductType,
		Meta:        req.Meta,
	}
}

type AdminPushNotificationRequest struct {
	UserID   int64  `form:"userId" validate:"required"`
	Title    string `form:"title" validate:"required"`
	Body     string `form:"body" validate:"required"`
	ImageURL string `form:"imageUrl"`
	Data     string `form:"data"`
}

type AdminCreateOrUpdateDestinationRequest struct {
	Name        string `form:"name" validate:"required"`
	Description string `form:"description"`
	Format      string `form:"format"`
	CheckerID   string `form:"checkerId"`
	Fields      []struct {
		Priority int    `form:"priority" validate:"required"`
		Key      string `form:"key" validate:"required"`
		Label    string `form:"label" validate:"required"`
		Type     string `form:"type" validate:"required"`
		Options  []struct {
			Label string `form:"label" validate:"required"`
			Value string `form:"value" validate:"required"`
		} `form:"options"`
		Description string `form:"description"`
		Required    bool   `form:"required"`
	} `form:"fields"`
}

func (req AdminCreateOrUpdateDestinationRequest) Model() model.CreateOrUpdateProductDestination {
	fields := map[string]model.DestinationField{}

	for _, field := range req.Fields {
		options := make([]model.DestinationFieldOption, 0)
		if field.Type == model.ProductDestinationFieldOptionType {
			for _, opt := range field.Options {
				options = append(options, model.DestinationFieldOption{
					Label: opt.Label,
					Value: opt.Value,
				})
			}
		}
		fields[field.Key] = model.DestinationField{
			Priority:    field.Priority,
			Label:       field.Label,
			Type:        field.Type,
			Description: field.Description,
			Options:     options,
			Required:    field.Required,
		}
	}

	return model.CreateOrUpdateProductDestination{
		Name:        req.Name,
		Description: &req.Description,
		Format:      &req.Format,
		CheckerID:   &req.CheckerID,
		Fields:      fields,
	}
}

type AdminCreateOrUpdatePulsaCategoryRequest struct {
	Prefix     string `form:"prefix" validate:"required,max=255"`
	CategoryID int64  `form:"categoryId" validate:"required"`
}

type AdminPurchasesRequest struct {
	StartDate        int64  `query:"startDate" validate:"required"`
	EndDate          int64  `query:"endDate" validate:"required"`
	DestinationQuery string `query:"destinationQuery"`
}

type AdminImportProductFromDigiflazzRequest struct {
	ToCategoryID  int64  `form:"toCategoryId" validate:"required"`
	ProductKind   string `form:"productKind"`
	DestinationID int64  `form:"destinationId" validate:"required"`
	AddPrice      int64  `form:"addPrice" validate:"required,min=1"`

	PrefixProductName string `form:"prefixProductName"`
	Category          string `form:"category" validate:"required"`
	Brand             string `form:"brand" validate:"required"`
	Type              string `form:"type" validate:"required"`
}

type AdminKycReviewRequest struct {
	ID       int64 `form:"id"`
	Accepted bool  `form:"accepted"`
}

type AdminReplaceMetaRequest struct {
	Name  string `form:"name"`
	Value string `form:"value"`
}
