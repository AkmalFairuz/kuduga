package paymentmethod

import (
	"github.com/akmalfairuz/finance/module/duitku"
	"time"
)

type BRIVA struct {
	client *duitku.Client
}

func NewBRIVA(client *duitku.Client) *BRIVA {
	return &BRIVA{client: client}
}

func (p *BRIVA) NoFee() bool {
	return false
}

func (p *BRIVA) ID() int {
	return BrivaID
}

func (p *BRIVA) Name() string {
	return "BRI Virtual Account"
}

func (p *BRIVA) ExpireDuration() time.Duration {
	return time.Hour * 24
}

func (p *BRIVA) CreatePayment(opt CreatePaymentOptions) (Data, error) {
	return createDuitkuPayment(p.client, p.ExpireDuration(), ToDuitkuPaymentCode(p.ID()), opt)
}

func (p *BRIVA) Check(refId string) (CheckResult, error) {
	return checkDuitkuPaymentStatus(p.client, refId)
}

func (p *BRIVA) PublicInfo(externalPaymentId, externalPaymentData string) (any, error) {
	return map[string]any{"paymentUrl": externalPaymentData}, nil
}

func (p *BRIVA) MinAmount() int64 {
	return 10000
}

func (p *BRIVA) MaxAmount() int64 {
	return 5000000
}

func (p *BRIVA) Description() string {
	return "Biaya admin sebesar Rp3.000. Status pembayaran secara otomatis akan dikonfirmasi setelah melakukan pembayaran bayar."
}

func (p *BRIVA) ImageUrl() string {
	return "https://s3.example.invalid/assets/briva.png"
}

func (p *BRIVA) Cancel(refId string) error {
	return nil
}
