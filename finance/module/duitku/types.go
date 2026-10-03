package duitku

const (
	SuccessCode   = "00"
	PendingCode   = "01"
	CancelledCode = "02"
)

const (
	SuccessResultCode = "00"
	FailedResultCode  = "01"
)

type CreatePaymentRequest struct {
	MerchantCode    string `json:"merchantCode"`
	PaymentAmount   int64  `json:"paymentAmount"`
	MerchantOrderID string `json:"merchantOrderId"`
	ProductDetails  string `json:"productDetails"`
	Email           string `json:"email"`
	PaymentMethod   string `json:"paymentMethod"`
	CustomerVaName  string `json:"customerVaName"`
	ReturnUrl       string `json:"returnUrl"`
	ExpiryPeriod    int64  `json:"expiryPeriod"`
	CallbackUrl     string `json:"callbackUrl"`
	Signature       string `json:"signature"`
}

type CreatePaymentResponse struct {
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`

	MerchantCode string `json:"merchantCode"`
	Reference    string `json:"reference"`
	PaymentUrl   string `json:"paymentUrl"`
	VaNumber     string `json:"vaNumber"`
	Amount       string `json:"amount"`
	QrString     string `json:"qrString"`
}

type CreatePaymentOptions struct {
	OrderID        string
	CustomerVaName string
	Email          string
	PaymentMethod  string
	ProductDetails string
	Amount         int64
	ExpiryPeriod   int64
}

type GetPaymentMethodRequest struct {
	MerchantCode string `json:"merchantCode"`
	Amount       int64  `json:"amount"`
	Datetime     string `json:"datetime"`
	Signature    string `json:"signature"`
}

type GetPaymentMethodResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`

	PaymentFee []struct {
		PaymentMethod string `json:"paymentMethod"`
		PaymentName   string `json:"paymentName"`
		PaymentImage  string `json:"paymentImage"`
		TotalFee      string `json:"totalFee"`
	}
}

type CheckStatusRequest struct {
	MerchantCode    string `json:"merchantCode"`
	MerchantOrderID string `json:"merchantOrderId"`
	Signature       string `json:"signature"`
}

type CheckStatusResponse struct {
	MerchantOrderID string `json:"merchantOrderId"`
	Reference       string `json:"reference"`
	Amount          int64  `json:"amount"`
	StatusCode      string `json:"statusCode"`
	StatusMessage   string `json:"statusMessage"`
}

type CallbackRequest struct {
	MerchantCode     string `form:"merchantCode"`
	Amount           int64  `form:"amount"`
	MerchantOrderID  string `form:"merchantOrderId"`
	ProductDetail    string `form:"productDetail"`
	AdditionalParam  string `form:"additionalParam"`
	PaymentCode      string `form:"paymentCode"`
	ResultCode       string `form:"resultCode"`
	MerchantUserID   string `form:"merchantUserId"`
	Reference        string `form:"reference"`
	Signature        string `form:"signature"`
	PublisherOrderID string `form:"publisherOrderId"`
	SpUserHash       string `form:"spUserHash"`
	SettlementDate   string `form:"settlementDate"`
	IssuerCode       string `form:"issuerCode"`
}
