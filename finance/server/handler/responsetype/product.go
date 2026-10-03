package responsetype

import (
	"github.com/akmalfairuz/finance/server/model"
	"sort"
)

type ProductResponse struct {
	ID                  int64  `json:"id"`
	CategoryID          int64  `json:"categoryId"`
	CategoryName        string `json:"categoryName"`
	CategoryParentID    *int64 `json:"categoryParentId"`
	Type                int    `json:"type"`
	ImageUrl            string `json:"imageUrl"`
	DestinationType     int64  `json:"destinationType"`
	Name                string `json:"name"`
	Kind                string `json:"kind"`
	Description         string `json:"description"`
	Price               int64  `json:"price"`
	BillAdmin           *int64 `json:"billAdmin,omitempty"`
	DestinationId       int64  `json:"destinationId"`
	BeforeDiscountPrice int64  `json:"beforeDiscountPrice,omitempty"`
	IsAvailable         bool   `json:"isAvailable"`
}

func NewProductResponseFromModel(product model.Product) ProductResponse {
	isAvailable := product.IsAvailable
	if product.IsMaintenance() {
		isAvailable = false
	}
	ret := ProductResponse{
		ID:                  product.ID,
		CategoryID:          product.CategoryID,
		CategoryName:        product.CategoryName,
		CategoryParentID:    product.CategoryParentID,
		Type:                product.Type,
		ImageUrl:            product.ImageUrl,
		DestinationType:     product.DestinationType,
		Name:                product.Name,
		Kind:                product.Kind,
		Description:         product.Description,
		Price:               product.Price,
		DestinationId:       product.DestinationType,
		BeforeDiscountPrice: product.BeforeDiscountPrice,
		IsAvailable:         isAvailable,
	}

	if product.IsBill() {
		ret.BillAdmin = &product.BillAdmin
	}

	return ret
}

type ProductCategoryResponse struct {
	ID          int64                     `json:"id"`
	ParentID    *int64                    `json:"parentId"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	IconUrl     string                    `json:"iconUrl"`
	ProductType int                       `json:"productType"`
	Children    []ProductCategoryResponse `json:"children"`
	Meta        map[string]string         `json:"meta"`
}

type ProductDestinationResponse struct {
	ID          int64                             `json:"id"`
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	Format      string                            `json:"format,omitempty"`
	CheckerID   *string                           `json:"checkerId"`
	Fields      []ProductDestinationFieldResponse `json:"fields"`
}

func NewProductDestinationResponseFromModel(destination model.ProductDestination, forAdmin bool) ProductDestinationResponse {
	format := ""
	if forAdmin {
		format = destination.Format_
	}
	fields := make([]ProductDestinationFieldResponse, 0)

	for key, field := range destination.Fields() {
		options := make([]ProductDestinationFieldOptionResponse, 0)
		for _, opt := range field.Options {
			options = append(options, ProductDestinationFieldOptionResponse{
				Label: opt.Label,
				Value: opt.Value,
			})
		}
		f := ProductDestinationFieldResponse{
			Priority:    field.Priority,
			Key:         key,
			Label:       field.Label,
			Type:        field.Type,
			Options:     options,
			Description: field.Description,
			Required:    field.Required,
		}
		fields = append(fields, f)
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Priority < fields[j].Priority
	})
	var checkerId *string
	if destination.CheckerID != "" {
		checkerId = &destination.CheckerID
	}
	return ProductDestinationResponse{
		ID:          destination.ID,
		Name:        destination.Name,
		Description: destination.Description,
		Format:      format,
		CheckerID:   checkerId,
		Fields:      fields,
	}
}

type ProductDestinationFieldResponse struct {
	Priority    int                                     `json:"priority,omitempty"`
	Key         string                                  `json:"key"`
	Label       string                                  `json:"label"`
	Type        string                                  `json:"type"`
	Description string                                  `json:"description"`
	Options     []ProductDestinationFieldOptionResponse `json:"options,omitempty"`
	Required    bool                                    `json:"required"`
}

type ProductDestinationFieldOptionResponse struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
