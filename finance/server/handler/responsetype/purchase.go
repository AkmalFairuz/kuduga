package responsetype

type CreatePurchaseResponse struct {
	PurchaseId int64 `json:"purchaseId"`
}

type PurchaseResponse struct {
	PurchaseId          int64  `json:"purchaseId"`
	ProductId           int64  `json:"productId"`
	ProductSku          string `json:"productSku"`
	ProductName         string `json:"productName"`
	ProductKind         string `json:"productKind"`
	ProductType         int    `json:"productType"`
	ProductDescription  string `json:"productDescription"`
	ProductCategoryId   int64  `json:"productCategoryId"`
	ProductCategoryName string `json:"productCategoryName"`
	Status              int    `json:"status"`
	StatusText          string `json:"statusText"`

	Destination map[string]string `json:"destination"`

	TotalBillFee int64               `json:"totalBillFee,omitempty"`
	TotalBill    int64               `json:"totalBill,omitempty"`
	BillAdmin    *int64              `json:"billAdmin,omitempty"`
	BillData     []map[string]string `json:"billData,omitempty"`

	ExtraData map[string]any `json:"extraData"`
	Note      string         `json:"note"`

	Proof            *string `json:"proof"`
	Price            int64   `json:"price"`
	UserSellPrice    int64   `json:"userSellPrice"`
	Fee              int64   `json:"fee"`
	PaymentMethod    *int    `json:"paymentMethod"`
	PaymentExpiredAt int64   `json:"paymentExpiredAt"`
	PaymentData      string  `json:"paymentData"`
	CreatedAt        int64   `json:"createdAt"`
}

type CheckBillResponse struct {
	ID           int64               `json:"id"`
	BillAmount   int64               `json:"billAmount"`
	Admin        int64               `json:"admin"`
	RealAdmin    int64               `json:"realAdmin"`
	TotalPrice   int64               `json:"totalPrice"`
	CustomerName string              `json:"customerName"`
	Data         map[string]string   `json:"data"`
	BillData     []map[string]string `json:"billData"`
}

type PublicPurchaseResponse struct {
	PurchaseId          int64  `json:"purchaseId"`
	ProductSku          string `json:"productSku"`
	ProductName         string `json:"productName"`
	ProductKind         string `json:"productKind"`
	ProductType         int    `json:"productType"`
	ProductCategoryName string `json:"productCategoryName"`
	Price               int64  `json:"price"`
	Status              int    `json:"status"`
	Destination         string `json:"destination"`
	CreatedAt           int64  `json:"createdAt"`
}
