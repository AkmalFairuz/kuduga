package provider

import (
	"github.com/akmalfairuz/finance/server/model"
	"github.com/valyala/fasthttp"
)

type SandboxPurchaseProvider struct {
}

func NewSandboxPurchaseProvider() *SandboxPurchaseProvider {
	return &SandboxPurchaseProvider{}
}

func (s SandboxPurchaseProvider) CheckBill(product model.Product, destination string, refId string) (*model.PurchaseCheckBillResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s SandboxPurchaseProvider) Process(product model.Product, purchase model.Purchase) (*model.PurchaseProcessResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s SandboxPurchaseProvider) CheckStatus(purchase model.Purchase) (*model.PurchaseCheckStatusResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s SandboxPurchaseProvider) HandleCallback(req *fasthttp.RequestCtx) (*model.PurchaseCheckStatusResult, error) {
	//TODO implement me
	panic("implement me")
}
