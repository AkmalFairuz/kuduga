package paymentmethod

import (
	"encoding/json"
	"fmt"
	"github.com/akmalfairuz/finance/module/fetch"
	"strings"
	"time"
)

type BankBCAConfig struct {
	ApiUrl         string `env:"BCA_PAYMENT_API_URL"`
	ApiKey         string `env:"BCA_PAYMENT_API_KEY"`
	CallbackApiKey string `env:"BCA_PAYMENT_CALLBACK_API_KEY"`
}

type BankBCA struct {
	config BankBCAConfig
}

func NewBankBCA(config BankBCAConfig) *BankBCA {
	return &BankBCA{
		config: config,
	}
}

func (p *BankBCA) NoFee() bool {
	return true
}

func (p *BankBCA) ID() int {
	return BankBcaID
}

func (p *BankBCA) Name() string {
	return "Bank BCA"
}

func (p *BankBCA) ExpireDuration() time.Duration {
	return time.Hour * 12
}

func (p *BankBCA) MinAmount() int64 {
	return 50000
}

func (p *BankBCA) MaxAmount() int64 {
	return 10000000
}

func (p *BankBCA) Description() string {
	return "Gratis biaya admin. Nominal yang harus di transfer akan bertambah sebesar 1 sampai 1000 (bukan biaya admin). Pembayaran akan di proses selama 5 menit setelah transfer dilakukan."
}

func (p *BankBCA) ImageUrl() string {
	return "https://upload.wikimedia.org/wikipedia/commons/5/5c/Bank_Central_Asia.svg"
}

// PublicInfo it shows the bank account destination
func (p *BankBCA) PublicInfo(_, externalPaymentData string) (any, error) {
	parse := strings.Split(externalPaymentData, " | ")
	if len(parse) != 2 {
		return nil, fmt.Errorf("invalid external payment data")
	}
	return map[string]any{
		"accountNumber": parse[0],
		"accountName":   parse[1],
	}, nil
}

func (p *BankBCA) CreatePayment(opt CreatePaymentOptions) (Data, error) {
	resp, err := fetch.PostJSON(p.config.ApiUrl+"/createPayment", map[string]any{
		"amount":  opt.Amount,
		"refId":   opt.RefID,
		"channel": "bca",
		"apiKey":  p.config.ApiKey,
	})
	if err != nil {
		return Data{}, err
	}

	if resp.Status != 200 {
		return Data{}, fmt.Errorf("got status code: %d, body: %s", resp.Status, string(resp.Body))
	}

	type BankBCAResponse struct {
		RefId         string `json:"refId"`
		Amount        int64  `json:"amount"`
		AccountNumber string `json:"accountNumber"`
		AccountName   string `json:"accountName"`
		ExpiredAt     int64  `json:"expiredAt"`
	}

	var bankBCAResponse BankBCAResponse
	if err := json.Unmarshal(resp.Body, &bankBCAResponse); err != nil {
		return Data{}, err
	}

	return Data{
		ExternalPaymentId:   bankBCAResponse.RefId,
		Amount:              bankBCAResponse.Amount,
		ExpiredAt:           bankBCAResponse.ExpiredAt,
		ExternalPaymentData: fmt.Sprintf("%s | %s", bankBCAResponse.AccountNumber, bankBCAResponse.AccountName),
	}, nil
}

func (p *BankBCA) Check(refId string) (CheckResult, error) {
	panic("implement me")
}

func (p *BankBCA) Cancel(refId string) error {
	resp, err := fetch.PostJSON(p.config.ApiUrl+"/cancelPayment", map[string]any{
		"refId":   refId,
		"channel": "bca",
		"apiKey":  p.config.ApiKey,
	})
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return fmt.Errorf("got status code: %d, body: %s", resp.Status, string(resp.Body))
	}

	return nil
}
