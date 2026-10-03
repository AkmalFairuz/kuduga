package provider

import (
	"github.com/akmalfairuz/finance/server/model"
	"github.com/valyala/fasthttp"
)

type PurchaseProvider interface {
	CheckBill(product model.Product, destination string, refId string) (*model.PurchaseCheckBillResult, error)
	Process(product model.Product, purchase model.Purchase) (*model.PurchaseProcessResult, error)
	CheckStatus(purchase model.Purchase) (*model.PurchaseCheckStatusResult, error)
	HandleCallback(req *fasthttp.RequestCtx) (*model.PurchaseCheckStatusResult, error)
}
