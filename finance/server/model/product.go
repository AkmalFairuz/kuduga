package model

import (
	"strconv"
	"strings"
	"time"
)

type ProductCategory struct {
	ID          int64  `db:"id"`
	ParentID    *int64 `db:"parentId"`
	Name        string `db:"name"`
	Description string `db:"description"`
	IconUrl     string `db:"iconUrl"`
	ProductType int    `db:"productType"`
	Meta_       string `db:"meta"`

	Children []ProductCategory `db:"-"`
}

func (p ProductCategory) Meta() map[string]string {
	m := map[string]string{}
	if p.Meta_ == "" {
		return m
	}
	for _, item := range strings.Split(p.Meta_, "\n") {
		x, y, _ := strings.Cut(item, "=")
		m[x] = y
	}
	return m
}

func (p ProductCategory) RequireKyc() bool {
	return p.Meta()["requireKyc"] == "true"
}

const (
	ProductTypeDigital         = 1
	ProductTypeDigitalPostpaid = 2
)

var ProductProviders = []string{"digiflazz", "tripay"}

type Product struct {
	ID int64 `db:"id"`

	CategoryID       int64  `db:"categoryId"`
	CategoryName     string `db:"categoryName"`
	CategoryParentID *int64 `db:"categoryParentId"`
	CategoryMeta     string `db:"categoryMeta"`

	Sku          string `db:"sku"`
	SkuPriority  int    `db:"skuPriority"`
	ExternalSku  string `db:"externalSku"`
	Provider     string `db:"provider"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	Type         int    `db:"type"`
	ImageUrl     string `db:"imageUrl"`
	Kind         string `db:"kind"`
	PurchaseNote string `db:"purchaseNote"`

	CutOffStart string `db:"cutOffStart"`
	CutOffEnd   string `db:"cutOffEnd"`

	DestinationType        int64  `db:"destinationType"`
	DestinationName        string `db:"destinationName"`
	DestinationDescription string `db:"destinationDescription"`

	ProofParser string `db:"proofParser"`

	IsAvailable bool `db:"isAvailable"`

	Price               int64 `db:"price"`
	WholesalePrice      int64 `db:"wholesalePrice"`
	MaxWholesalePrice   int64 `db:"maxWholesalePrice"`
	BeforeDiscountPrice int64 `db:"beforeDiscountPrice"`

	// Only available for bill product
	BillAdmin int64 `db:"billAdmin"`

	CreatedAt int64 `db:"createdAt"`
	UpdatedAt int64 `db:"updatedAt"`
}

func (p Product) IsBill() bool {
	return p.Type == ProductTypeDigitalPostpaid
}

func (p Product) IsMaintenance() bool {
	return isInHHMM(time.Now().UTC().Add(time.Hour*7), p.CutOffStart, p.CutOffEnd)
}

func isInHHMM(t time.Time, startHHMM, endHHMM string) bool {
	if startHHMM == endHHMM || startHHMM == "" || endHHMM == "" {
		return false
	}

	hh, mm, _ := t.Clock()

	currentM := (hh * 60) + mm

	startM := convHHMM(startHHMM)
	endM := convHHMM(endHHMM)

	if startM > endM {
		return currentM >= startM || currentM <= endM
	}

	return currentM >= startM && currentM <= endM
}

func convHHMM(hhMM string) int {
	x, y, _ := strings.Cut(hhMM, ":")
	x2, err := strconv.Atoi(x)
	if err != nil {
		panic(err)
	}
	y2, err := strconv.Atoi(y)
	if err != nil {
		panic(err)
	}
	return (x2 * 60) + y2
}

func (p Product) Category() ProductCategory {
	return ProductCategory{
		ID:       p.CategoryID,
		ParentID: p.CategoryParentID,
		Name:     p.CategoryName,
		Meta_:    p.CategoryMeta,
	}
}

func (p Product) StatusCode() string {
	if p.IsMaintenance() {
		return "CUT"
	}
	if !p.IsAvailable {
		return "OFF"
	}
	return "ON"
}

type CreateOrUpdateProduct struct {
	CategoryID          int64
	Sku                 string
	SkuPriority         *int
	ExternalSku         string
	Provider            string
	Name                string
	Kind                *string
	PurchaseNote        *string
	Description         *string
	Type                int
	CutOffStart         *string
	CutOffEnd           *string
	ImageUrl            *string
	DestinationType     int64
	IsAvailable         *bool
	Price               int64
	BillAdmin           *int64
	MaxWholesalePrice   int64
	WholesalePrice      int64
	BeforeDiscountPrice *int64
	ProofParser         *string
}

type CreateOrUpdateProductCategory struct {
	Name        string
	ParentID    **int64
	Description *string
	IconUrl     *string
	ProductType int
	Meta        *string
}
