package paymentmethod

import (
	"github.com/akmalfairuz/finance/module/duitku"
	"github.com/akmalfairuz/finance/module/helper"
	"time"
)

var duitkuPaymentCodes = map[int]string{
	AlfamartID:           "FT",
	BrivaID:              "BR",
	BniVirtualID:         "I1",
	PermataBankVirtualID: "BT",
	AtmBersamaID:         "A1",
	QrisID:               "SP",
}

func ToDuitkuPaymentCode(paymentMethod int) string {
	return duitkuPaymentCodes[paymentMethod]
}

func FromDuitkuCode(paymentCode string) int {
	for k, v := range duitkuPaymentCodes {
		if v == paymentCode {
			return k
		}
	}
	return 0
}

func createDuitkuPayment(duitkuClient *duitku.Client, expireIn time.Duration, paymentMethod string, opt CreatePaymentOptions) (Data, error) {
	fee, err := duitkuClient.GetFee(paymentMethod, opt.Amount)
	if err != nil {
		return Data{}, err
	}
	result, err := duitkuClient.CreatePayment(duitku.CreatePaymentOptions{
		OrderID:        opt.RefID,
		CustomerVaName: opt.CustomerName,
		PaymentMethod:  paymentMethod,
		ProductDetails: opt.ProductDetails,
		Amount:         opt.Amount,
		ExpiryPeriod:   expireIn.Milliseconds() / 1000 / 60,
	})
	if err != nil {
		return Data{}, err
	}
	return Data{
		ExternalPaymentId:   opt.RefID, // TODO: fix this
		ExternalPaymentData: result.PaymentUrl,
		Amount:              int64(helper.StringToInt(result.Amount)),
		Fee:                 fee,
		ExpiredAt:           time.Now().Add(expireIn).Unix(),
	}, nil
}

func checkDuitkuPaymentStatus(duitkuClient *duitku.Client, orderId string) (CheckResult, error) {
	result, err := duitkuClient.CheckStatus(orderId)
	if err != nil {
		return CheckResult{}, err
	}
	ret := CheckResult{}
	switch result.StatusCode {
	case duitku.SuccessCode:
		ret.Status = StatusSuccess
	case duitku.PendingCode:
		ret.Status = StatusWaitingPayment
	case duitku.CancelledCode:
		ret.Status = StatusFailed
	}
	return ret, nil
}
