package service

import (
	"fmt"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/akmalfairuz/finance/server/repository"
)

type NotificationService struct {
	notificationRepository *repository.NotificationRepository
	notificationProvider   provider.NotificationProvider
	authService            *AuthService
}

func NewNotificationService(notificationRepository *repository.NotificationRepository, notificationProvider provider.NotificationProvider) *NotificationService {
	return &NotificationService{notificationRepository: notificationRepository, notificationProvider: notificationProvider}
}

func (s *NotificationService) SetAuthService(authService *AuthService) {
	s.authService = authService
}

func (s *NotificationService) GetNotificationsByUserId(userId int64) ([]model.Notification, error) {
	return s.notificationRepository.Get(userId)
}

func (s *NotificationService) MarkAsReadNotification(userId int64, id int64) error {
	return s.notificationRepository.MarkRead(userId, id)
}

func (s *NotificationService) BroadcastNotification(opt *provider.PushNotificationOptions) error {
	return s.PushNotificationToTopic("global", opt)
}

func (s *NotificationService) PushNotificationToTopic(topic string, opt *provider.PushNotificationOptions) error {
	return s.notificationProvider.PushNotificationToTopic(topic, opt)
}

func (s *NotificationService) CreateNotification(create model.CreateNotification) error {
	notificationId, err := s.notificationRepository.Create(create)
	if err != nil {
		return err
	}
	if create.Push {
		if create.Data == nil {
			create.Data = map[string]string{}
		}
		create.Data["_notificationId"] = fmt.Sprintf("%d", notificationId)
		if err := s.PushNotificationToUser(create.UserID, &provider.PushNotificationOptions{
			Notification: &provider.NotificationOptions{
				Title: create.Title,
				Body:  create.Description,
			},
			Data: create.Data,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *NotificationService) PushNotificationToUser(userId int64, opt *provider.PushNotificationOptions) error {
	tokens, err := s.authService.GetTokens(userId)
	if err != nil {
		return err
	}

	sentDeviceTokens := map[string]bool{}

	for _, token := range tokens {
		if token.DeviceFcmToken != nil {
			if _, ok := sentDeviceTokens[*token.DeviceFcmToken]; ok {
				continue
			}
			sentDeviceTokens[*token.DeviceFcmToken] = true
			if err := s.notificationProvider.PushNotificationToDevice(*token.DeviceFcmToken, opt); err != nil {
				if err.Error() == "Requested entity was not found." {
					if err := s.authService.RevokeToken(token.ID); err != nil {
						return err
					}
					continue
				}
				return err
			}
		}
	}

	return nil
}

func (s *NotificationService) InitializeFcmToken(fcmToken string) error {
	return s.notificationProvider.SubscribeToTopic(fcmToken, "global")
}
