package service

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/random"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/akmalfairuz/finance/server/repository"
)

type UserKycService struct {
	userKycRepository *repository.UserKycRepository
	storage           provider.StorageProvider

	userService         *UserService
	notificationService *NotificationService
	alert               provider.TextAlertProvider
}

var userKycBucketName = "user-kyc"

func NewUserKycService(userKycRepository *repository.UserKycRepository, storage provider.StorageProvider) *UserKycService {
	return &UserKycService{userKycRepository: userKycRepository, storage: storage}
}

func (s *UserKycService) SetNotificationService(notificationService *NotificationService) {
	s.notificationService = notificationService
}

func (s *UserKycService) SetUserService(userService *UserService) {
	s.userService = userService
}

func (s *UserKycService) SetAlert(alert provider.TextAlertProvider) {
	s.alert = alert
}

func (s *UserKycService) Create(create *model.CreateUserKyc) error {
	create.DocumentFileID = fmt.Sprintf("%d-%s.png", create.UserID, random.String(80))
	kycId, err := s.userKycRepository.Create(create)
	if err != nil {
		return err
	}
	if err := s.storage.Put(userKycBucketName, create.DocumentFileID, create.DocumentFile); err != nil {
		return err
	}
	if s.alert != nil {
		_ = s.alert.Alert(fmt.Sprintf("[New Kyc Request #%d]\nUserID: %d\nFull Name: %s", kycId, create.UserID, create.FullName))
	}
	return nil
}

func (s *UserKycService) UpdateStatus(kyc model.UserKyc, status int) error {
	if err := s.userKycRepository.UpdateStatus(kyc.ID, status); err != nil {
		return err
	}

	if err := s.userService.UpdateUserKycStatus(kyc.UserID, status == model.UserKycStatusSuccess); err != nil {
		return err
	}
	return nil
}

func (s *UserKycService) GetByUserID(userId int64) (model.UserKyc, error) {
	return s.userKycRepository.GetByUserId(userId)
}

func (s *UserKycService) Get(id int64) (model.UserKyc, error) {
	return s.userKycRepository.Get(id)
}

func (s *UserKycService) GetAllPending() ([]model.UserKyc, error) {
	return s.userKycRepository.GetAllPending()
}

func (s *UserKycService) GetAttachmentPrefix() string {
	return s.storage.PublicEndpoint() + "/" + userKycBucketName
}
