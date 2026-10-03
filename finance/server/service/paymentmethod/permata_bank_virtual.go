package paymentmethod

import (
	"github.com/akmalfairuz/finance/module/duitku"
	"time"
)

type PermataBankVirtual struct {
	client *duitku.Client
}

func NewPermataBankVirtual(client *duitku.Client) *PermataBankVirtual {
	return &PermataBankVirtual{client: client}
}

func (p *PermataBankVirtual) NoFee() bool {
	return false
}

func (p *PermataBankVirtual) ID() int {
	return PermataBankVirtualID
}

func (p *PermataBankVirtual) Name() string {
	return "Permata Bank Virtual Account"
}

func (p *PermataBankVirtual) ExpireDuration() time.Duration {
	return time.Hour * 24
}

func (p *PermataBankVirtual) CreatePayment(opt CreatePaymentOptions) (Data, error) {
	return createDuitkuPayment(p.client, p.ExpireDuration(), ToDuitkuPaymentCode(p.ID()), opt)
}

func (p *PermataBankVirtual) Check(refId string) (CheckResult, error) {
	return checkDuitkuPaymentStatus(p.client, refId)
}

func (p *PermataBankVirtual) PublicInfo(externalPaymentId, externalPaymentData string) (any, error) {
	return map[string]any{"paymentUrl": externalPaymentData}, nil
}

func (p *PermataBankVirtual) MinAmount() int64 {
	return 10000
}

func (p *PermataBankVirtual) MaxAmount() int64 {
	return 5000000
}

func (p *PermataBankVirtual) Description() string {
	return "Biaya admin sebesar Rp3.000. Status pembayaran secara otomatis akan dikonfirmasi setelah melakukan pembayaran bayar."
}

func (p *PermataBankVirtual) ImageUrl() string {
	return "https://s3.example.invalid/assets/permata_bank.png"
}

func (p *PermataBankVirtual) Cancel(refId string) error {
	return nil
}
