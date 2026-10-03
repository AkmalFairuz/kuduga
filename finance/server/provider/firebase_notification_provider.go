package provider

import (
	"context"
	"firebase.google.com/go/v4/messaging"
)

type FirebaseNotificationProvider struct {
	client *messaging.Client
}

func NewFirebaseNotificationProvider(client *messaging.Client) *FirebaseNotificationProvider {
	return &FirebaseNotificationProvider{client: client}
}

func (p *FirebaseNotificationProvider) toFirebaseMessage(opt *PushNotificationOptions) *messaging.Message {
	var notification *messaging.Notification
	if opt.Notification != nil {
		notification = &messaging.Notification{
			Title:    opt.Notification.Title,
			Body:     opt.Notification.Body,
			ImageURL: opt.Notification.ImageURL,
		}
	}

	message := &messaging.Message{
		Data:         opt.Data,
		Notification: notification,
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				ChannelID:  "other", // TODO: implement
				Priority:   messaging.PriorityHigh,
				Visibility: messaging.VisibilityPublic,
			},
		},
	}
	return message
}

func (p *FirebaseNotificationProvider) PushNotificationToTopic(topic string, opt *PushNotificationOptions) error {
	message := p.toFirebaseMessage(opt)
	message.Topic = topic

	if _, err := p.client.Send(context.Background(), message); err != nil {
		return err
	}

	return nil
}

func (p *FirebaseNotificationProvider) PushNotificationToDevice(deviceToken string, opt *PushNotificationOptions) error {
	message := p.toFirebaseMessage(opt)
	message.Token = deviceToken

	if _, err := p.client.Send(context.Background(), message); err != nil {
		return err
	}

	return nil
}

func (p *FirebaseNotificationProvider) SubscribeToTopic(deviceToken string, topic string) error {
	_, err := p.client.SubscribeToTopic(context.TODO(), []string{deviceToken}, topic)
	return err
}

func (p *FirebaseNotificationProvider) UnsubscribeToTopic(deviceToken string, topic string) error {
	_, err := p.client.UnsubscribeFromTopic(context.TODO(), []string{deviceToken}, topic)
	return err
}
