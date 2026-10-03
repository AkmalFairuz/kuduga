package discordwebhook

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/fetch"
)

type Client struct {
	WebhookUrl string
}

func New(webhookUrl string) *Client {
	return &Client{
		WebhookUrl: webhookUrl,
	}
}

func (c *Client) SendText(message string) error {
	return c.Send(Message{
		Content: message,
	})
}

func (c *Client) Send(message Message) error {
	res, err := fetch.PostJSON(c.WebhookUrl, message)
	if err != nil {
		return err
	}
	if res.IsStatusError() {
		return fmt.Errorf("discord webhook error [%d]: %s", res.Status, res.Body)
	}
	return nil
}
