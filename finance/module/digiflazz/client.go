package digiflazz

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/akmalfairuz/finance/module/fetch"
	"github.com/akmalfairuz/finance/module/hash"
	"github.com/valyala/fasthttp"
	"strings"
)

type Config struct {
	Testing       bool   `env:"DIGIFLAZZ_TESTING" envDefault:"true"`
	ApiKey        string `env:"DIGIFLAZZ_API_KEY"`
	Secret        string `env:"DIGIFLAZZ_SECRET"`
	Username      string `env:"DIGIFLAZZ_USERNAME"`
	Proxy         string `env:"DIGIFLAZZ_PROXY"`
	IP            string `env:"DIGIFLAZZ_IP"`
	CallbackUrl   string `env:"DIGIFLAZZ_CB_URL"`
	CallbackToken string `env:"DIGIFLAZZ_CB_TOKEN"`
}

const EndpointV1 = "https://api.digiflazz.com/v1"

type Client struct {
	cfg Config
}

func New(cfg Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) GetBalance() (int64, error) {
	var resp CheckBalanceResponse

	req := &CheckBalanceRequest{
		Cmd:      "deposit",
		Username: c.cfg.Username,
		Sign:     hash.Md5(c.cfg.Username + c.cfg.ApiKey + "depo"),
	}

	if _, err := c.req("/cek-saldo", &req, &resp); err != nil {
		return 0, err
	}

	return resp.Deposit, nil
}

func (c *Client) GetProducts() ([]ProductBaseResponse, error) {
	ret := make([]ProductBaseResponse, 0)

	prepaid, err := c.GetPrepaidPriceList("")
	if err != nil {
		return nil, err
	}
	for _, p := range prepaid {
		ret = append(ret, p)
	}

	postpaid, err := c.GetPostpaidPriceList("")
	if err != nil {
		return nil, err
	}
	for _, p := range postpaid {
		ret = append(ret, p)
	}

	return ret, nil
}

func (c *Client) GetPrepaidPriceList(code string) ([]ProductResponse, error) {
	var resp []ProductResponse
	cmd := "prepaid"
	if _, err := c.req("/price-list", &PriceListRequest{
		Cmd:      cmd,
		Username: c.cfg.Username,
		Code:     code,
		Sign:     hash.Md5(c.cfg.Username + c.cfg.ApiKey + cmd),
	}, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetPostpaidPriceList(code string) ([]ProductPostpaidResponse, error) {
	var resp []ProductPostpaidResponse
	cmd := "pasca"
	if _, err := c.req("/price-list", &PriceListRequest{
		Cmd:      cmd,
		Username: c.cfg.Username,
		Code:     code,
		Sign:     hash.Md5(c.cfg.Username + c.cfg.ApiKey + cmd),
	}, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) Deposit(amount int64, bank string, ownerName string) (DepositResponse, error) {
	var resp DepositResponse

	type req struct {
		Username  string `json:"username"`
		Amount    int64  `json:"amount"`
		Bank      string `json:"Bank"`
		OwnerName string `json:"owner_name"`
		Sign      string `json:"sign"`
	}

	if _, err := c.req("/deposit", &req{
		Username:  c.cfg.Username,
		Amount:    amount,
		Bank:      bank,
		OwnerName: ownerName,
		Sign:      hash.Md5(c.cfg.Username + c.cfg.ApiKey + "deposit"),
	}, &resp); err != nil {
		return DepositResponse{}, err
	}

	return resp, nil
}

func (c *Client) Topup(request TopupRequest) (TopupResponse, error) {
	request.Username = c.cfg.Username
	request.Sign = hash.Md5(c.cfg.Username + c.cfg.ApiKey + request.RefId)
	request.CbUrl = c.cfg.CallbackUrl + "?token=" + c.cfg.CallbackToken
	request.Testing = c.cfg.Testing

	var resp TopupResponse
	if _, err := c.req("/transaction", request, &resp); err != nil {
		return TopupResponse{}, err
	}

	return resp, nil
}

func (c *Client) CheckStatus(request TopupRequest) (TopupResponse, error) {
	request.MaxPrice = 1
	return c.Topup(request)
}

func (c *Client) CheckBill(request BillRequest) (BillResponse, error) {
	return c.reqBill(request, "inq-pasca")
}

func (c *Client) PayBill(request BillRequest) (BillResponse, error) {
	return c.reqBill(request, "pay-pasca")
}

func (c *Client) reqBill(request BillRequest, commands string) (BillResponse, error) {
	request.Commands = commands
	request.Username = c.cfg.Username
	request.Sign = hash.Md5(c.cfg.Username + c.cfg.ApiKey + request.RefId)
	request.Testing = c.cfg.Testing

	var resp BillResponse
	if _, err := c.req("/transaction", request, &resp); err != nil {
		return BillResponse{}, err
	}

	return resp, nil
}

func (c *Client) InquiryPLN(customerNo string) (InquiryPLNResponse, error) {
	request := InquiryPLNRequest{
		CustomerNo: customerNo,
		Commands:   "pln-subscribe",
	}

	var resp InquiryPLNResponse
	if _, err := c.req("/transaction", request, &resp); err != nil {
		return InquiryPLNResponse{}, err
	}

	return resp, nil
}

func (c *Client) HandleCallback(ctx *fasthttp.RequestCtx) (WebhookEvent, error) {
	//TODO: check ip
	//if c.cfg.IP != "" {
	//	if string(ctx.Request.Header.Peek("X-Real-IP")) != c.cfg.IP {
	//		return WebhookEvent{}, errors.New("invalid ip")
	//	}
	//}

	if string(ctx.QueryArgs().Peek("token")) != c.cfg.CallbackToken {
		return WebhookEvent{}, errors.New("invalid callback token")
	}

	req := &ctx.Request
	signature := string(req.Header.Peek("X-Hub-Signature"))
	userAgent := string(req.Header.Peek("User-Agent"))
	event := string(req.Header.Peek("X-Digiflazz-Event"))

	if signature == "" || userAgent == "" || event == "" {
		return WebhookEvent{}, errors.New("invalid headers")
	}

	expectedSignature := "sha1=" + hash.HmacSha1(c.cfg.Secret, string(req.Body()))

	if expectedSignature != signature {
		fmt.Printf("signature: %s, got %s\n", expectedSignature, signature)
		return WebhookEvent{}, errors.New("invalid signature")
	}

	if !strings.HasPrefix(strings.ToLower(userAgent), "digiflazz") {
		return WebhookEvent{}, errors.New("invalid user agent")
	}

	var wresp wrappedResponse
	if err := json.Unmarshal(req.Body(), &wresp); err != nil {
		return WebhookEvent{}, err
	}

	return WebhookEvent{
		Event: event,
		Data:  wresp.Data,
	}, nil
}

func (c *Client) DecodeStatusUpdateEvent(data json.RawMessage) (StatusUpdateEvent, error) {
	var event StatusUpdateEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return StatusUpdateEvent{}, err
	}
	return event, nil
}

func (c *Client) req(path string, body any, dst any) (*fetch.Response, error) {
	resp, err := fetch.PostJSON(EndpointV1+path, body)
	if err != nil {
		return nil, err
	}

	var wresp wrappedResponse
	if err := json.Unmarshal(resp.Body, &wresp); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(wresp.Data, dst); err != nil {
		return nil, err
	}

	return resp, nil
}
