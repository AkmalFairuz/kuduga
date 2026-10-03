package digiflazz

type CheckBalanceRequest struct {
	Cmd      string `json:"cmd"`
	Username string `json:"username"`
	Sign     string `json:"sign"`
}

type PriceListRequest struct {
	Cmd      string `json:"cmd"`
	Username string `json:"username"`
	Code     string `json:"code,omitempty"`
	Sign     string `json:"sign"`
}

type TopupRequest struct {
	Username     string `json:"username"`
	BuyerSkuCode string `json:"buyer_sku_code"`
	CustomerNo   string `json:"customer_no"`
	RefId        string `json:"ref_id"`
	// Sign md5(username + api_key + ref_id)
	Sign     string `json:"sign"`
	Testing  bool   `json:"testing,omitempty"`
	MaxPrice int64  `json:"max_price,omitempty"`
	CbUrl    string `json:"cb_url"`
}

type BillRequest struct {
	Commands     string `json:"commands"`
	Username     string `json:"username"`
	BuyerSkuCode string `json:"buyer_sku_code"`
	CustomerNo   string `json:"customer_no"`
	RefId        string `json:"ref_id"`
	Sign         string `json:"sign"`
	Testing      bool   `json:"testing,omitempty"`
}

type InquiryPLNRequest struct {
	// CustomerNo or nomor meteran
	CustomerNo string `json:"customer_no"`
	// Commands value must be pln-subscribe
	Commands string `json:"commands"`
}
