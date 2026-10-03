package duitku

type Config struct {
	ApiKey        string `env:"DUITKU_API_KEY"`
	MerchantCode  string `env:"DUITKU_MERCHANT_CODE"`
	CallbackUrl   string `env:"DUITKU_CALLBACK_URL"`
	CallbackToken string `env:"DUITKU_CALLBACK_TOKEN"`
	ReturnUrl     string `env:"DUITKU_RETURN_URL"`
	Mode          string `env:"DUITKU_MODE"`
}
