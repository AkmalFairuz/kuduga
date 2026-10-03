package paymentmethod

import (
	"github.com/akmalfairuz/finance/module/duitku"
	"time"
)

type ATMBersama struct {
	client *duitku.Client
}

func NewATMBersama(client *duitku.Client) *ATMBersama {
	return &ATMBersama{client: client}
}

func (p *ATMBersama) NoFee() bool {
	return false
}

func (p *ATMBersama) ID() int {
	return AtmBersamaID
}

func (p *ATMBersama) Name() string {
	return "Virtual Account Bank Lain"
}

func (p *ATMBersama) ExpireDuration() time.Duration {
	return time.Hour * 24
}

func (p *ATMBersama) CreatePayment(opt CreatePaymentOptions) (Data, error) {
	return createDuitkuPayment(p.client, p.ExpireDuration(), ToDuitkuPaymentCode(p.ID()), opt)
}

func (p *ATMBersama) Check(refId string) (CheckResult, error) {
	return checkDuitkuPaymentStatus(p.client, refId)
}

func (p *ATMBersama) PublicInfo(externalPaymentId, externalPaymentData string) (any, error) {
	return map[string]any{"paymentUrl": externalPaymentData}, nil
}

func (p *ATMBersama) MinAmount() int64 {
	return 10000
}

func (p *ATMBersama) MaxAmount() int64 {
	return 5000000
}

func (p *ATMBersama) Description() string {
	return "Biaya admin sebesar Rp3.000. Status pembayaran secara otomatis akan dikonfirmasi setelah melakukan pembayaran bayar."
}

func (p *ATMBersama) ImageUrl() string {
	return "https://s3.example.invalid/assets/atmbersama.png"
}

func (p *ATMBersama) Cancel(refId string) error {
	return nil
}
