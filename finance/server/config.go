package server

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/module/digiflazz"
	"github.com/akmalfairuz/finance/module/duitku"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/akmalfairuz/finance/server/service/paymentmethod"
)

type Config struct {
	AppMode string `env:"APP_MODE" envDefault:"development"`

	ListenAddress    string   `env:"LISTEN_ADDRESS" envDefault:":3001"`
	TrustedProxies   []string `env:"TRUSTED_PROXIES" envSeparator:","`
	CORSAllowOrigins string   `env:"CORS_ALLOW_ORIGINS" envDefault:"http://localhost:3000"`

	Database database.Config

	Redis struct {
		Address  string `env:"REDIS_ADDRESS"`
		Password string `env:"REDIS_PASSWORD"`
	}

	Minio               provider.MinioConfig
	AWSCredentials      provider.AWSCredentials
	BcaPayment          paymentmethod.BankBCAConfig
	Digiflazz           digiflazz.Config
	OTPEmailFrom        string `env:"OTP_EMAIL_FROM"`
	Duitku              duitku.Config
	GoogleOAuthClientId string `env:"GOOGLE_OAUTH_CLIENT_ID"`

	DiscordWebhook DiscordWebhookConfig
}

func (c Config) IsProduction() bool {
	return c.AppMode == "production"
}

type DiscordWebhookConfig struct {
	General           string `env:"DISCORD_WEBHOOK_GENERAL"`
	Deposit           string `env:"DISCORD_WEBHOOK_DEPOSIT"`
	Purchase          string `env:"DISCORD_WEBHOOK_PURCHASE"`
	ProductManagement string `env:"DISCORD_WEBHOOK_PRODUCT_MANAGEMENT"`
}
