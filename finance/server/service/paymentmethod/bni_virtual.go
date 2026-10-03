package paymentmethod

import (
	"github.com/akmalfairuz/finance/module/duitku"
	"time"
)

type BNIVirtual struct {
	client *duitku.Client
}

func NewBNIVirtual(client *duitku.Client) *BNIVirtual {
	return &BNIVirtual{client: client}
}

func (p *BNIVirtual) NoFee() bool {
	return false
}

func (p *BNIVirtual) ID() int {
	return BniVirtualID
}

func (p *BNIVirtual) Name() string {
	return "BNI Virtual Account"
}

func (p *BNIVirtual) ExpireDuration() time.Duration {
	return time.Hour * 24
}

func (p *BNIVirtual) CreatePayment(opt CreatePaymentOptions) (Data, error) {
	return createDuitkuPayment(p.client, p.ExpireDuration(), ToDuitkuPaymentCode(p.ID()), opt)
}

func (p *BNIVirtual) Check(refId string) (CheckResult, error) {
	return checkDuitkuPaymentStatus(p.client, refId)
}

func (p *BNIVirtual) PublicInfo(externalPaymentId, externalPaymentData string) (any, error) {
	return map[string]any{"paymentUrl": externalPaymentData}, nil
}

func (p *BNIVirtual) MinAmount() int64 {
	return 10000
}

func (p *BNIVirtual) MaxAmount() int64 {
	return 5000000
}

func (p *BNIVirtual) Description() string {
	return "Biaya admin sebesar Rp3.000. Status pembayaran secara otomatis akan dikonfirmasi setelah melakukan pembayaran bayar."
}

func (p *BNIVirtual) ImageUrl() string {
	return "https://s3.example.invalid/assets/bni.png"
}

func (p *BNIVirtual) Cancel(refId string) error {
	return nil
}
