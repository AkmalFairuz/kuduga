package provider

import "github.com/akmalfairuz/finance/module/discordwebhook"

type DiscordTextAlertProvider struct {
	client *discordwebhook.Client
}

func NewDiscordTextAlertProvider(webhookUrl string) *DiscordTextAlertProvider {
	return &DiscordTextAlertProvider{client: discordwebhook.New(webhookUrl)}
}

func (p *DiscordTextAlertProvider) Alert(text string) error {
	return p.client.SendText(text)
}
