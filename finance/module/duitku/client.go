package duitku

import (
	"encoding/json"
	"fmt"
	"github.com/akmalfairuz/finance/module/fetch"
	"github.com/akmalfairuz/finance/module/hash"
	"github.com/akmalfairuz/finance/module/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	sandboxUrl    = "https://sandbox.duitku.com/webapi/api/merchant"
	productionUrl = "https://passport.duitku.com/webapi/api/merchant"
)

type Client struct {
	log    logrus.FieldLogger
	config Config
}

func NewClient(log logrus.FieldLogger, config Config) *Client {
	return &Client{log: log.WithField("module", "duitku"), config: config}
}

func (c *Client) GetFee(paymentMethod string, amount int64) (int64, error) {
	resp, err := c.GetPaymentMethod(amount)
	if err != nil {
		return 0, err
	}
	for _, p := range resp.PaymentFee {
		if p.PaymentMethod == paymentMethod {
			return int64(helper.StringToInt(p.TotalFee)), nil
		}
	}
	return 0, fmt.Errorf("payment method %s not found", paymentMethod)
}

func (c *Client) GetPaymentMethod(amount int64) (GetPaymentMethodResponse, error) {
	var dst GetPaymentMethodResponse
	dt := time.Now().Format("2006-01-02 15:04:05")
	resp, err := c.req(&dst, "POST", "/paymentmethod/getpaymentmethod", GetPaymentMethodRequest{
		MerchantCode: c.config.MerchantCode,
		Amount:       amount,
		Datetime:     dt,
		Signature:    hash.Sha256(fmt.Sprintf("%s%d%s%s", c.config.MerchantCode, amount, dt, c.config.ApiKey)),
	})
	if err != nil {
		return GetPaymentMethodResponse{}, err
	}
	if resp.Status != http.StatusOK {
		if dst.ResponseMessage == "" {
			return GetPaymentMethodResponse{}, fmt.Errorf("duitku error: %d", resp.Status)
		}
		return GetPaymentMethodResponse{}, fmt.Errorf("duitku error: %s - %s", dst.ResponseCode, dst.ResponseMessage)
	}
	return dst, nil
}

func (c *Client) CreatePayment(options CreatePaymentOptions) (CreatePaymentResponse, error) {
	var dst CreatePaymentResponse
	resp, err := c.req(&dst, "POST", "/v2/inquiry", &CreatePaymentRequest{
		MerchantCode:    c.config.MerchantCode,
		PaymentAmount:   options.Amount,
		MerchantOrderID: options.OrderID,
		ProductDetails:  options.ProductDetails,
		Email:           options.Email,
		PaymentMethod:   options.PaymentMethod,
		CustomerVaName:  options.CustomerVaName,
		ReturnUrl:       c.config.ReturnUrl,
		ExpiryPeriod:    options.ExpiryPeriod,
		CallbackUrl:     c.config.CallbackUrl + "?token=" + c.config.CallbackToken,
		Signature:       hash.Md5(fmt.Sprintf("%s%s%d%s", c.config.MerchantCode, options.OrderID, options.Amount, c.config.ApiKey)),
	})
	if err != nil {
		c.log.Errorf("failed to create payment: %+v", err)
		return CreatePaymentResponse{}, err
	}

	if resp.Status != http.StatusOK {
		if dst.StatusMessage == "" || dst.StatusCode == "" {
			return CreatePaymentResponse{}, fmt.Errorf("duitku unknown error: %d", resp.Status)
		}
		return CreatePaymentResponse{}, fmt.Errorf("duitku error: %s - %s", dst.StatusCode, dst.StatusMessage)
	}

	c.log.WithFields(logrus.Fields{
		"orderId":   options.OrderID,
		"reference": dst.Reference,
		"amount":    options.Amount,
		"payment":   options.PaymentMethod,
	}).Infof("payment created")

	return dst, nil
}

func (c *Client) CheckStatus(orderId string) (CheckStatusResponse, error) {
	var dst CheckStatusResponse
	resp, err := c.req(&dst, "POST", "/transactionStatus", &CheckStatusRequest{
		MerchantCode:    c.config.MerchantCode,
		MerchantOrderID: orderId,
		Signature:       hash.Md5(fmt.Sprintf("%s%s%s", c.config.MerchantCode, orderId, c.config.ApiKey)),
	})
	if err != nil {
		return CheckStatusResponse{}, err
	}
	if resp.Status != http.StatusOK {
		if dst.StatusMessage == "" || dst.StatusCode == "" {
			return CheckStatusResponse{}, fmt.Errorf("duitku unknown error: %d", resp.Status)
		}
		return CheckStatusResponse{}, fmt.Errorf("duitku error: %s - %s", dst.StatusCode, dst.StatusMessage)
	}
	return dst, nil
}

func (c *Client) HandleCallback(ctx *fiber.Ctx) (CallbackRequest, error) {
	// validate token (internal)
	if ctx.Query("token") != c.config.CallbackToken {
		c.log.Errorf("invalid token, expected: %s, got: %s", c.config.CallbackToken, ctx.Query("token"))
		return CallbackRequest{}, fmt.Errorf("invalid token")
	}

	var dst CallbackRequest
	if err := ctx.BodyParser(&dst); err != nil {
		c.log.Errorf("failed to parse callback body: %+v", err)
		return CallbackRequest{}, err
	}

	log := c.log.WithFields(logrus.Fields{
		"reference":   dst.Reference,
		"orderId":     dst.MerchantOrderID,
		"amount":      dst.Amount,
		"paymentCode": dst.PaymentCode,
	})

	if dst.MerchantCode != c.config.MerchantCode {
		return CallbackRequest{}, fmt.Errorf("merchant code mismatch")
	}
	expectedSign := hash.Md5(fmt.Sprintf("%s%d%s%s", c.config.MerchantCode, dst.Amount, dst.MerchantOrderID, c.config.ApiKey))
	if dst.Signature != expectedSign {
		return CallbackRequest{}, fmt.Errorf("signature mismatch")
	}

	log.Infof("received callback from duitku resultCode: %s", dst.ResultCode)

	return dst, nil
}

func (c *Client) req(dst any, method string, path string, body any) (*fetch.Response, error) {
	url := ""
	if c.config.Mode == "production" {
		url = productionUrl
	} else {
		url = sandboxUrl
	}

	opt := &fetch.Options{
		Headers: fetch.Map{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Method: method,
		Url:    fmt.Sprintf("%s%s", url, path),
	}

	if body != nil {
		bytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		opt.Body = bytes
	}

	resp, err := fetch.Request(opt)
	if err != nil {
		return nil, err
	}

	if dst != nil {
		if err = json.Unmarshal(resp.Body, dst); err != nil {
			return nil, err
		}
	}

	return resp, nil
}
