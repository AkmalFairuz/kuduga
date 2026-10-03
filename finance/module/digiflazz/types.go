package digiflazz

import "encoding/json"

type wrappedResponse struct {
	Data json.RawMessage `json:"data"`
}

type StatusUpdateEvent struct {
	TrxId          string `json:"trx_id"`
	RefId          string `json:"ref_id"`
	BuyerSkuCode   string `json:"buyer_sku_code"`
	Rc             string `json:"rc"`
	Sn             string `json:"sn"`
	Status         string `json:"status"`
	Message        string `json:"message"`
	CustomerNo     string `json:"customer_no"`
	BuyerLastSaldo int64  `json:"buyer_last_saldo"`
	Price          int64  `json:"price"`
}
