package digiflazz

type CheckBalanceResponse struct {
	Deposit int64 `json:"deposit"`
}

type ProductBaseResponse interface {
	GetProductName() string
	GetCategory() string
	GetBrand() string
	GetSellerName() string
	GetBuyerSkuCode() string
	GetBuyerProductStatus() bool
	GetSellerProductStatus() bool
	GetDesc() string
	GetPrice() int64
}

type ProductResponse struct {
	ProductName         string `json:"product_name"`
	Category            string `json:"category"`
	Brand               string `json:"brand"`
	Type                string `json:"type"`
	SellerName          string `json:"seller_name"`
	Price               int64  `json:"price"`
	BuyerSkuCode        string `json:"buyer_sku_code"`
	BuyerProductStatus  bool   `json:"buyer_product_status"`
	SellerProductStatus bool   `json:"seller_product_status"`
	UnlimitedStock      bool   `json:"unlimited_stock"`
	Stock               int64  `json:"stock"`
	Multi               bool   `json:"multi"`
	StartCutOff         string `json:"start_cut_off"`
	EndCutOff           string `json:"end_cut_off"`
	Desc                string `json:"desc"`
}

func (p ProductResponse) GetProductName() string {
	return p.ProductName
}

func (p ProductResponse) GetCategory() string {
	return p.Category
}

func (p ProductResponse) GetBrand() string {
	return p.Brand
}

func (p ProductResponse) GetSellerName() string {
	return p.SellerName
}

func (p ProductResponse) GetBuyerSkuCode() string {
	return p.BuyerSkuCode
}

func (p ProductResponse) GetBuyerProductStatus() bool {
	return p.BuyerProductStatus
}

func (p ProductResponse) GetSellerProductStatus() bool {
	return p.SellerProductStatus
}

func (p ProductResponse) GetDesc() string {
	return p.Desc
}

func (p ProductResponse) GetPrice() int64 {
	return p.Price
}

type ProductPostpaidResponse struct {
	ProductName         string `json:"product_name"`
	Category            string `json:"category"`
	Brand               string `json:"brand"`
	SellerName          string `json:"seller_name"`
	Admin               int64  `json:"admin"`
	Commission          int64  `json:"commission"`
	BuyerSkuCode        string `json:"buyer_sku_code"`
	BuyerProductStatus  bool   `json:"buyer_product_status"`
	SellerProductStatus bool   `json:"seller_product_status"`
	Desc                string `json:"desc"`
}

func (p ProductPostpaidResponse) GetProductName() string {
	return p.ProductName
}

func (p ProductPostpaidResponse) GetCategory() string {
	return p.Category
}

func (p ProductPostpaidResponse) GetBrand() string {
	return p.Brand
}

func (p ProductPostpaidResponse) GetSellerName() string {
	return p.SellerName
}

func (p ProductPostpaidResponse) GetBuyerSkuCode() string {
	return p.BuyerSkuCode
}

func (p ProductPostpaidResponse) GetBuyerProductStatus() bool {
	return p.BuyerProductStatus
}

func (p ProductPostpaidResponse) GetSellerProductStatus() bool {
	return p.SellerProductStatus
}

func (p ProductPostpaidResponse) GetDesc() string {
	return p.Desc
}

func (p ProductPostpaidResponse) GetPrice() int64 {
	return p.Admin - p.Commission
}

type DepositResponse struct {
	Rc     string `json:"rc"`
	Amount int64  `json:"amount"`
	Notes  string `json:"notes"`
}

type TopupResponse struct {
	RefId        string `json:"ref_id"`
	CustomerNo   string `json:"customer_no"`
	BuyerSkuCode string `json:"buyer_sku_code"`
	Message      string `json:"message"`
	// Status = Sukses, Pending, Gagal
	Status         string `json:"status"`
	Rc             string `json:"rc"`
	Sn             string `json:"sn"`
	BuyerLastSaldo int64  `json:"buyer_last_saldo"`
	Price          int64  `json:"price"`
	Tele           string `json:"tele"`
	Wa             string `json:"wa"`
}

type BillResponse struct {
	RefId          string `json:"ref_id"`
	CustomerNo     string `json:"customer_no"`
	CustomerName   string `json:"customer_name"`
	BuyerSkuCode   string `json:"buyer_sku_code"`
	Admin          int64  `json:"admin"`
	Message        string `json:"message"`
	Status         string `json:"status"`
	Rc             string `json:"rc"`
	BuyerLastSaldo int64  `json:"buyer_last_saldo"`
	Price          int64  `json:"price"`
	SellingPrice   int64  `json:"selling_price"`

	// Desc only available if request commands is inq-pasca
	Desc map[string]any `json:"desc"`

	// Sn only available if request commands is pay-pasca
	Sn string `json:"sn"`
	// Detail only available if request commands is pay-pasca
	Detail map[string]any `json:"detail"`
}

type InquiryPLNResponse struct {
	CustomerNo   string `json:"customer_no"`
	Name         string `json:"name"`
	SegmentPower string `json:"segment_power"`
	SubscriberId string `json:"subscriber_id"`
	MeterNo      string `json:"meter_no"`
}
