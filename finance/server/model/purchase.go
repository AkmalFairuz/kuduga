package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	PurchaseStatusWaitingPayment = 0
	PurchaseStatusProcess        = 1
	PurchaseStatusSuccess        = 2
	PurchaseStatusFailed         = 3
)

type Purchase struct {
	ID int64 `db:"id"`
	// UserID is the user who made the purchase
	// If it is null, it means the purchase is made by guest
	UserID *int64 `db:"userId"`
	// RefID is the reference ID of the purchase from third party
	RefID string `db:"refId"`
	// ContactEmail is the email of the user who made the purchase
	ContactEmail        string `db:"contactEmail"`
	ProductID           int64  `db:"productId"`
	ProductName         string `db:"productName"`
	ProductKind         string `db:"productKind"`
	ProductDescription  string `db:"productDescription"`
	ProductType         int    `db:"productType"`
	ProductSku          string `db:"productSku"`
	ProductExternalSku  string `db:"productExternalSku"`
	ProductProvider     string `db:"productProvider"`
	ProductCategoryId   int64  `db:"productCategoryId"`
	ProductCategoryName string `db:"productCategoryName"`
	ProviderInfo        string `db:"providerInfo"`
	Price               int64  `db:"price"`
	UserSellPrice       int64  `db:"userSellPrice"`

	// TotalBill is total bill exclude bill fee. Only available if product is postpaid
	TotalBill int64 `db:"totalBill"`
	// TotalBillFee is total bill fee. Only available if product is postpaid
	TotalBillFee int64  `db:"totalBillFee"`
	BillAdmin    int64  `db:"billAdmin"`
	BillData_    string `db:"billData"`

	Fee            int64 `db:"fee"`
	WholesalePrice int64 `db:"wholesalePrice"`

	// Destination is the destination of the purchase
	// for example, if the product is a mobile credit, the destination is the phone number
	Destination          string `db:"destination"`
	DestinationId        int64  `db:"destinationId"`
	DetailedDestination_ string `db:"detailedDestination"`
	ReadableDestination_ string `db:"readableDestination"`

	// Proof is the proof of the purchase
	Proof *string `db:"proof"`
	// ParsedProof is the parsed proof of the purchase
	// for example, if the product is a mobile credit, the parsed proof is the serial number
	ParsedProof string `db:"parsedProof"`

	// PaymentMethod is the payment method used to make the purchase
	// It's null if the payment method using user balance
	PaymentMethod *int `db:"paymentMethod"`
	// PaymentRefId is the reference ID of the payment from third party
	// If user pay with balance, the payment ref id contains the transaction id
	PaymentRefId string `db:"paymentRefId"`
	// PaymentData is the data of the payment from third party
	PaymentData      string `db:"paymentData"`
	PaymentExpiredAt int64  `db:"paymentExpiredAt"`

	Refunded   bool   `db:"refunded"`
	ExtraData_ string `db:"extraData"`
	Note       string `db:"note"`

	// InternalData_ contains internal information
	// usually contains supplier customer service, etc.
	InternalData_ string `db:"internalData"`

	Status int `db:"status"`

	SuccessAt int64 `db:"successAt"`
	CreatedAt int64 `db:"createdAt"`
}

func (p Purchase) DetailedDestination() (map[string]string, error) {
	if p.DetailedDestination_ == "" {
		return map[string]string{}, nil
	}
	var detail map[string]string
	if err := json.Unmarshal([]byte(p.DetailedDestination_), &detail); err != nil {
		return nil, err
	}
	return detail, nil
}

func (p Purchase) ReadableDestination() (map[string]string, error) {
	if p.ReadableDestination_ == "" {
		return map[string]string{}, nil
	}
	var detail map[string]string
	if err := json.Unmarshal([]byte(p.ReadableDestination_), &detail); err != nil {
		return nil, err
	}
	return detail, nil
}

func (p Purchase) ExtraData() (map[string]any, error) {
	if p.ExtraData_ == "" {
		return map[string]any{}, nil
	}
	var extraData map[string]any
	if err := json.Unmarshal([]byte(p.ExtraData_), &extraData); err != nil {
		return nil, err
	}
	return extraData, nil
}

func (p Purchase) BillData() ([]map[string]string, error) {
	if p.BillData_ == "" {
		return []map[string]string{}, nil
	}
	var ret []map[string]string
	if err := json.Unmarshal([]byte(p.BillData_), &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (p Purchase) StatusText() string {
	return PurchaseStatusText(p.Status)
}

func PurchaseStatusText(status int) string {
	switch status {
	case PurchaseStatusWaitingPayment:
		return "Menunggu Pembayaran"
	case PurchaseStatusProcess:
		return "Proses"
	case PurchaseStatusSuccess:
		return "Sukses"
	case PurchaseStatusFailed:
		return "Gagal"
	}
	return "-"
}

type CreatePurchase struct {
	Product              Product
	ProductCategory      ProductCategory
	ProviderInfo         string
	UserID               *int64
	RefId                string
	ContactEmail         string
	Price                int64
	UserSellPrice        int64
	TotalBill            int64
	TotalBillFee         int64
	BillAdmin            int64
	Fee                  int64
	WholesalePrice       int64
	Destination          string
	DetailedDestination_ map[string]string
	ReadableDestination_ map[string]string
	PaymentMethod        *int
	PaymentRefId         string
	PaymentData          string
	PaymentExpiredAt     int64
	BillData_            []map[string]string
	ExtraData_           map[string]any
	InternalData_        map[string]any
}

func (create CreatePurchase) DetailedDestination() (string, error) {
	if len(create.DetailedDestination_) == 0 {
		return "", nil
	}
	bytes, err := json.Marshal(create.DetailedDestination_)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (create CreatePurchase) ReadableDestination() (string, error) {
	if len(create.ReadableDestination_) == 0 {
		return "", nil
	}
	bytes, err := json.Marshal(create.ReadableDestination_)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (create CreatePurchase) ExtraData() (string, error) {
	if len(create.ExtraData_) == 0 {
		return "", nil
	}
	bytes, err := json.Marshal(create.ExtraData_)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (create CreatePurchase) BillData() (string, error) {
	bytes, err := json.Marshal(create.BillData_)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (create CreatePurchase) InternalData() (string, error) {
	bytes, err := json.Marshal(create.InternalData_)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type PurchaseTrack struct {
	ID          int64  `db:"id"`
	PurchaseID  int64  `db:"purchaseId"`
	NewStatus   *int   `db:"newStatus"`
	Description string `db:"description"`
}

type CreatePurchaseOptions struct {
	UserID          *int64
	UserKycVerified bool
	Email           string
	Destination     map[string]string
	PaymentMethod   *int
	ProductID       int64
	Bill            *BillPrePurchase
	AllowDuplicate  bool
	MaxPrice        int64
}

type CreatePurchaseResult struct {
	PurchaseID       int64
	Price            int64
	RefId            string
	PaymentExpiredAt int64
}

const (
	PurchaseProcessFailedReasonUnknown            = iota + 1
	PurchaseProcessFailedReasonInvalidDestination = iota + 1
)

type PurchaseProcessResult struct {
	Status  int
	Message string
	// Price is the wholesale price of the product
	// It must not show the price to the user
	Price int64
	// Proof is the proof of the purchase
	// for example, if the product is a mobile credit, the proof is serial number
	Proof        string
	FailedReason int
}

type PurchaseCheckStatusResult struct {
	RefId   string
	Status  int
	Message string
	Proof   string
}

type PurchasePaymentReceivedResult struct {
	Proof string
}

type GetUserPurchasesOptions struct {
	Status   *int
	BeforeID int64
	Limit    int64
}

type BillPrePurchase struct {
	ID           int64  `db:"id"`
	RefID        string `db:"refId"`
	UserID       *int64 `db:"userId"`
	ContactEmail string `db:"contactEmail"`

	ProductID int64 `db:"productId"`

	Destination          string `db:"destination"`
	DestinationID        int64  `db:"destinationId"`
	DetailedDestination_ string `db:"detailedDestination"`
	ReadableDestination_ string `db:"readableDestination"`

	BillAmount     int64 `db:"billAmount"`
	BillAdmin      int64 `db:"billAdmin"`
	Price          int64 `db:"price"`
	WholesalePrice int64 `db:"wholesalePrice"`

	BillData_  string `db:"billData"`
	ExtraData_ string `db:"extraData"`
	CreatedAt  int64  `db:"createdAt"`
	ExpiredAt  int64  `db:"expiredAt"`
}

func (model BillPrePurchase) BillData() ([]map[string]string, error) {
	var ret []map[string]string
	if err := json.Unmarshal([]byte(model.BillData_), &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (model BillPrePurchase) ExtraData() (map[string]any, error) {
	var ret map[string]any
	if err := json.Unmarshal([]byte(model.ExtraData_), &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type CreateBillPrePurchase struct {
	RefID                string
	UserID               *int64
	ContactEmail         string
	ProductID            int64
	Destination          string
	DestinationID        int64
	DetailedDestination_ map[string]string
	ReadableDestination_ map[string]string
	BillAmount           int64
	BillAdmin            int64
	Price                int64
	WholesalePrice       int64
	BillData_            []map[string]string
	ExtraData_           map[string]any
	ExpiredAt            int64
}

func (create CreateBillPrePurchase) DetailedDestination() (string, error) {
	if create.DetailedDestination_ == nil || len(create.DetailedDestination_) == 0 {
		return "", nil
	}
	bytes, err := json.Marshal(create.DetailedDestination_)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (create CreateBillPrePurchase) ReadableDestination() (string, error) {
	if create.ReadableDestination_ == nil || len(create.ReadableDestination_) == 0 {
		return "", nil
	}
	bytes, err := json.Marshal(create.ReadableDestination_)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (create CreateBillPrePurchase) ExtraData() (string, error) {
	if create.ExtraData_ == nil || len(create.ExtraData_) == 0 {
		return "", nil
	}
	bytes, err := json.Marshal(create.ExtraData_)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (create CreateBillPrePurchase) BillData() (string, error) {
	if create.BillData_ == nil || len(create.BillData_) == 0 {
		return "", nil
	}
	bytes, err := json.Marshal(create.BillData_)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type PurchaseCheckBillResult struct {
	ID            int64
	BillAmount    int64
	Price         int64
	Admin         int64
	RealAdmin     int64
	CustomerPrice int64
	CustomerName  string
	Data          map[string]string
	BillData      []map[string]string
}

func (model PurchaseCheckBillResult) AsExtraData() map[string]any {
	ret := map[string]any{}
	ret["Nama"] = model.CustomerName
	for k, v := range model.Data {
		ret[k] = v
	}
	return ret
}

func (model PurchaseCheckBillResult) RawData() string {
	ret := ""
	for k, v := range model.Data {
		ret += fmt.Sprintf("%s = %s\n", k, v)
	}
	for i, data := range model.BillData {
		ret += fmt.Sprintf("---------- Tagihan #%d ----------\n", i+1)
		for k, v := range data {
			ret += fmt.Sprintf("%s = %s\n", k, v)
		}
		ret += "---------------------------------\n"
	}
	return strings.TrimRight(ret, "\n")
}

type SearchPurchases struct {
	FromDate         int64
	ToDate           int64
	DestinationQuery string
}

type UserPurchaseStats struct {
	Count  int64 `db:"count"`
	Total  int64 `db:"total"`
	Profit int64 `db:"profit"`
}
