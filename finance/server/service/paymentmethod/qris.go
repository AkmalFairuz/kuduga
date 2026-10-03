package paymentmethod

import (
	"github.com/akmalfairuz/finance/module/duitku"
	"github.com/akmalfairuz/finance/module/helper"
	"time"
)

type Qris struct {
	duitkuClient *duitku.Client
}

func NewQris(duitkuClient *duitku.Client) *Qris {
	return &Qris{
		duitkuClient: duitkuClient,
	}
}

func (p *Qris) NoFee() bool {
	return false
}

func (p *Qris) ID() int {
	return QrisID
}

func (p *Qris) Name() string {
	return "QRIS"
}

func (p *Qris) ExpireDuration() time.Duration {
	return time.Hour
}

func (p *Qris) PublicInfo(_, externalPaymentData string) (any, error) {
	return map[string]string{
		"qrString": externalPaymentData,
	}, nil
}

func (p *Qris) CreatePayment(opt CreatePaymentOptions) (Data, error) {
	pmId := ToDuitkuPaymentCode(p.ID())

	// TODO: qris fee bug!
	//fee, err := p.duitkuClient.GetFee(pmId, opt.Amount)
	//if err != nil {
	//	return Data{}, err
	//}

	// estimate fee
	fee := int64(float64(opt.Amount)*0.007) + 1
	// TODO: THIS IS BAD! FIX THIS WHEN DUITKU GIVE US THE FEE VIA API

	result, err := p.duitkuClient.CreatePayment(duitku.CreatePaymentOptions{
		OrderID:        opt.RefID,
		CustomerVaName: opt.CustomerName,
		PaymentMethod:  pmId,
		ProductDetails: opt.ProductDetails,
		Amount:         opt.Amount,
		ExpiryPeriod:   p.ExpireDuration().Milliseconds() / 1000 / 60,
	})
	if err != nil {
		return Data{}, err
	}

	return Data{
		ExternalPaymentId:   opt.RefID, // TODO: fix this
		ExternalPaymentData: result.QrString,
		Amount:              int64(helper.StringToInt(result.Amount)),
		Fee:                 fee,
		ExpiredAt:           time.Now().Add(p.ExpireDuration()).Unix(),
	}, nil
}

func (p *Qris) MinAmount() int64 {
	return 10000
}

func (p *Qris) MaxAmount() int64 {
	return 5000000
}

func (p *Qris) Description() string {
	return "Biaya admin sebesar 0.7% dari nominal deposit"
}

func (p *Qris) ImageUrl() string {
	return "https://s3.example.invalid/assets/qris.svg"
}

func (p *Qris) Check(refId string) (CheckResult, error) {
	status, err := checkDuitkuPaymentStatus(p.duitkuClient, refId)
	if err != nil {
		return CheckResult{}, err
	}
	return status, nil
}

func (p *Qris) Cancel(refId string) error {
	return nil
}
