package provider

import (
	"errors"
	"fmt"
	"github.com/akmalfairuz/finance/module/digiflazz"
	"github.com/akmalfairuz/finance/module/helper"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"unicode"
)

var DigiflazzPurchaseProviderName = "digiflazz"

type DigiflazzPurchaseProvider struct {
	log    logrus.FieldLogger
	client *digiflazz.Client
}

func NewDigiflazzPurchaseProvider(log logrus.FieldLogger, client *digiflazz.Client) *DigiflazzPurchaseProvider {
	return &DigiflazzPurchaseProvider{
		log:    log,
		client: client,
	}
}

func (p *DigiflazzPurchaseProvider) Process(product model.Product, purchase model.Purchase) (*model.PurchaseProcessResult, error) {
	if product.Type == model.ProductTypeDigitalPostpaid {
		return p.payBill(product, purchase)
	}
	topup, err := p.client.Topup(digiflazz.TopupRequest{
		BuyerSkuCode: product.ExternalSku,
		CustomerNo:   purchase.Destination,
		RefId:        purchase.RefID,
		MaxPrice:     product.WholesalePrice,
	})
	if err != nil {
		return nil, err
	}

	ret := &model.PurchaseProcessResult{
		Price:   topup.Price,
		Message: topup.Message,
	}

	switch topup.Status {
	case "Sukses":
		ret.Status = model.PurchaseStatusSuccess
		ret.Proof = topup.Sn
		return ret, nil
	case "Gagal":
		ret.Status = model.PurchaseStatusFailed

		p.log.WithFields(logrus.Fields{
			"purchaseId":  purchase.ID,
			"externalSku": product.ExternalSku,
		}).Errorf("topup digiflazz failed: %+v", topup.Rc)

		switch topup.Rc {
		case "54":
			ret.FailedReason = model.PurchaseProcessFailedReasonInvalidDestination
		default:
			ret.FailedReason = model.PurchaseProcessFailedReasonUnknown
		}
		return ret, nil
	case "Pending":
		ret.Status = model.PurchaseStatusProcess
		return ret, nil
	default:
		return nil, fmt.Errorf("invalid status: %s", topup.Status)
	}
}

func (p *DigiflazzPurchaseProvider) payBill(product model.Product, purchase model.Purchase) (*model.PurchaseProcessResult, error) {
	bill, err := p.client.PayBill(digiflazz.BillRequest{
		BuyerSkuCode: product.ExternalSku,
		CustomerNo:   purchase.Destination,
		RefId:        purchase.RefID,
	})
	if err != nil {
		return nil, err
	}
	ret := &model.PurchaseProcessResult{
		Price:   bill.Price,
		Message: bill.Message,
	}
	switch bill.Status {
	case "Sukses":
		ret.Status = model.PurchaseStatusSuccess
		ret.Proof = bill.Sn
		return ret, nil
	case "Gagal":
		ret.Status = model.PurchaseStatusFailed
		return ret, nil
	case "Pending":
		ret.Status = model.PurchaseStatusProcess
		return ret, nil
	default:
		return nil, fmt.Errorf("invalid status: %s", bill.Status)
	}
}

func (p *DigiflazzPurchaseProvider) HandleCallback(ctx *fasthttp.RequestCtx) (*model.PurchaseCheckStatusResult, error) {
	res, err := p.client.HandleCallback(ctx)
	if err != nil {
		return nil, err
	}

	switch res.Event {
	case "update":
		result, err := p.client.DecodeStatusUpdateEvent(res.Data)
		if err != nil {
			return nil, err
		}
		var status int
		switch result.Status {
		case "Sukses":
			status = model.PurchaseStatusSuccess
		case "Gagal":
			status = model.PurchaseStatusFailed
		case "Pending":
			status = model.PurchaseStatusProcess
		default:
			return nil, fmt.Errorf("invalid status: %s", result.Status)
		}
		return &model.PurchaseCheckStatusResult{
			RefId:   result.RefId,
			Message: result.Message,
			Status:  status,
			Proof:   result.Sn,
		}, nil
	default:
		return nil, nil
	}
}

func (p *DigiflazzPurchaseProvider) CheckBill(product model.Product, destination string, refId string) (*model.PurchaseCheckBillResult, error) {
	res, err := p.client.CheckBill(digiflazz.BillRequest{
		BuyerSkuCode: product.ExternalSku,
		CustomerNo:   destination,
		RefId:        refId,
	})
	if err != nil {
		return nil, err
	}
	if res.Status != "Sukses" {
		switch res.Rc {
		case "60":
			return nil, errortype.ErrPurchaseBillNotFound
		case "54":
			return nil, errortype.ErrPurchaseInvalidDestination
		}
		return nil, fmt.Errorf("failed to check bill: status=%v rc=%v message=%v", res.Status, res.Rc, res.Message)
	}
	billAmount := res.SellingPrice - res.Admin
	lembarTagihan := int64(res.Desc["lembar_tagihan"].(float64))

	billData := make([]map[string]string, 0)
	detail := res.Desc["detail"].([]any)
	realAdmin := int64(0)
	for _, detail2 := range detail {
		data := map[string]string{}
		for k, v := range detail2.(map[string]any) {
			vStr := fmt.Sprintf("%v", v)
			if k == "admin" {
				vInt, _ := strconv.Atoi(vStr)
				realAdmin += int64(vInt)
				continue
			}
			formatRupiah := false
			if k == "denda" {
				if vStr == "0" {
					continue
				}
				formatRupiah = true
			}
			if k == "nilai_tagihan" {
				formatRupiah = true
			}
			if formatRupiah {
				vInt, _ := strconv.Atoi(vStr)
				if vInt > 0 {
					vStr = helper.FormatRupiah(int64(vInt))
				}
			}
			data[convertToTitleCase(k)] = vStr
		}
		billData = append(billData, data)
	}

	data := map[string]string{}
	for k, v := range res.Desc {
		if k == "lembar_tagihan" || k == "detail" {
			continue
		}
		data[convertToTitleCase(k)] = fmt.Sprintf("%v", v)
	}

	totalAdmin := product.Price * lembarTagihan
	customerPrice := billAmount + totalAdmin
	if customerPrice < res.Price {
		return nil, errors.New("check bill error: customer price less than digiflazz price")
	}

	return &model.PurchaseCheckBillResult{
		BillAmount:    billAmount,
		Price:         res.Price,
		Admin:         totalAdmin,
		RealAdmin:     realAdmin,
		CustomerName:  res.CustomerName,
		CustomerPrice: customerPrice,
		Data:          data,
		BillData:      billData,
	}, nil
}

func (p *DigiflazzPurchaseProvider) CheckStatus(purchase model.Purchase) (*model.PurchaseCheckStatusResult, error) {
	res, err := p.client.CheckStatus(digiflazz.TopupRequest{
		BuyerSkuCode: purchase.ProductExternalSku,
		CustomerNo:   purchase.Destination,
		RefId:        purchase.RefID,
		MaxPrice:     2,
	})
	if err != nil {
		return nil, err
	}
	var status int
	switch res.Status {
	case "Sukses":
		status = model.PurchaseStatusSuccess
	case "Gagal":
		status = model.PurchaseStatusFailed
	case "Pending":
		status = model.PurchaseStatusProcess
	default:
		return nil, fmt.Errorf("invalid status: %s", res.Status)
	}
	return &model.PurchaseCheckStatusResult{
		RefId:   res.RefId,
		Message: res.Message,
		Status:  status,
		Proof:   res.Sn,
	}, nil
}

func convertToTitleCase(input string) string {
	words := strings.Split(input, "_")
	for i, word := range words {
		r := []rune(word)
		r[0] = unicode.ToUpper(r[0])
		words[i] = string(r)
	}
	result := strings.Join(words, " ")

	return result
}
