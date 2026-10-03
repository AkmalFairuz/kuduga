package provider

type NotificationOptions struct {
	Title    string
	Body     string
	ImageURL string
}

type PushNotificationOptions struct {
	Notification *NotificationOptions
	Data         map[string]string
}

type NotificationProvider interface {
	PushNotificationToDevice(deviceToken string, opt *PushNotificationOptions) error
	PushNotificationToTopic(topic string, opt *PushNotificationOptions) error
	SubscribeToTopic(deviceToken string, topic string) error
	UnsubscribeToTopic(deviceToken string, topic string) error
}
