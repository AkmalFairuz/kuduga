package paymentmethod

import (
	"github.com/akmalfairuz/finance/module/duitku"
	"time"
)

type Alfamart struct {
	client *duitku.Client
}

func NewAlfamart(client *duitku.Client) *Alfamart {
	return &Alfamart{client: client}
}

func (p *Alfamart) NoFee() bool {
	return false
}

func (p *Alfamart) ID() int {
	return AlfamartID
}

func (p *Alfamart) Name() string {
	return "Alfamart"
}

func (p *Alfamart) ExpireDuration() time.Duration {
	return time.Hour * 23
}

func (p *Alfamart) CreatePayment(opt CreatePaymentOptions) (Data, error) {
	return createDuitkuPayment(p.client, p.ExpireDuration(), ToDuitkuPaymentCode(p.ID()), opt)
}

func (p *Alfamart) Check(refId string) (CheckResult, error) {
	return checkDuitkuPaymentStatus(p.client, refId)
}

func (p *Alfamart) PublicInfo(externalPaymentId, externalPaymentData string) (any, error) {
	return map[string]any{"paymentUrl": externalPaymentData}, nil
}

func (p *Alfamart) MinAmount() int64 {
	return 10000
}

func (p *Alfamart) MaxAmount() int64 {
	return 2000000
}

func (p *Alfamart) Description() string {
	return "Biaya admin sebesar Rp2.500"
}

func (p *Alfamart) ImageUrl() string {
	return "https://s3.example.invalid/assets/alfamart.png"
}

func (p *Alfamart) Cancel(refId string) error {
	return nil
}
