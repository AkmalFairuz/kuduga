package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/akmalfairuz/finance/module/random"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/akmalfairuz/finance/server/repository"
	"github.com/redis/go-redis/v9"
	"strings"
)

type SupportService struct {
	supportTicketRepository *repository.SupportTicketRepository
	fileStorage             provider.StorageProvider

	redisClient *redis.Client
	alert       provider.TextAlertProvider
}

var supportTicketBucketName = "support-ticket"

func NewSupportService(supportTicketRepository *repository.SupportTicketRepository, fileStorage provider.StorageProvider, redisClient *redis.Client) *SupportService {
	return &SupportService{supportTicketRepository: supportTicketRepository, fileStorage: fileStorage, redisClient: redisClient}
}

func (s *SupportService) SetAlert(alert provider.TextAlertProvider) {
	s.alert = alert
}

func (s *SupportService) GetAttachmentPrefix() string {
	return s.fileStorage.PublicEndpoint() + "/" + supportTicketBucketName
}

func (s *SupportService) processAttachment(ticketId int64, create *model.CreateSupportTicketMessage) error {
	if len(create.Attachments) != 0 {
		for _, attachment := range create.Attachments {
			path := fmt.Sprintf("%d__%s/%s", ticketId, random.String(64), attachment.Name)
			if err := s.fileStorage.Put(supportTicketBucketName, path, attachment.Bytes); err != nil {
				return err
			}
			create.RawAttachment_ += path + "\n"
		}
		create.RawAttachment_ = strings.TrimRight(create.RawAttachment_, "\n")
		create.Attachments = nil
	}
	return nil
}

func (s *SupportService) CreateSupportTicket(create model.CreateSupportTicket, initialMessage *model.CreateSupportTicketMessage) (int64, error) {
	ticketId, err := s.supportTicketRepository.CreateSupportTicket(create)
	if err != nil {
		return 0, err
	}
	initialMessage.TicketID = ticketId
	if err := s.processAttachment(ticketId, initialMessage); err != nil {
		return 0, err
	}
	if _, err := s.supportTicketRepository.CreateSupportTicketMessage(initialMessage); err != nil {
		return 0, err
	}
	_ = s.alert.Alert(fmt.Sprintf("[New support ticket #%d]\nAuthor: %s\n\n%s", ticketId, initialMessage.Author, initialMessage.Message))
	return ticketId, nil
}

func (s *SupportService) GetSupportTicketMessages(ticketId int64) ([]model.SupportTicketMessage, error) {
	return s.supportTicketRepository.GetSupportTicketMessages(ticketId)
}

func (s *SupportService) GetSupportTicketMessagesAfterID(ticketId int64, afterId int64) ([]model.SupportTicketMessage, error) {
	return s.supportTicketRepository.GetSupportTicketMessagesAfterID(ticketId, afterId)
}

func (s *SupportService) GetSupportTicketsByUserID(userId int64) ([]model.SupportTicket, error) {
	return s.supportTicketRepository.GetSupportTicketsByUserID(userId)
}

func (s *SupportService) GetSupportTicketByID(id int64) (model.SupportTicket, error) {
	return s.supportTicketRepository.GetSupportTicketByID(id)
}

func (s *SupportService) CreateSupportTicketMessage(create *model.CreateSupportTicketMessage) (int64, error) {
	ticket, err := s.supportTicketRepository.GetSupportTicketByID(create.TicketID)
	if err != nil {
		return 0, err
	}
	if ticket.Status == model.SupportTicketCloseStatus {
		return 0, errors.New("ticket is closed")
	}
	if err := s.processAttachment(ticket.ID, create); err != nil {
		return 0, err
	}

	_ = s.NotifySupportTicketUpdate(ticket.ID) // Notify admin
	_ = s.alert.Alert(fmt.Sprintf("[Support Ticket Update #%d]\nNew Message Alert from %s", create.TicketID, create.Author))

	return s.supportTicketRepository.CreateSupportTicketMessage(create)
}

func (s *SupportService) GetSupportTicketMessage(messageId int64) (model.SupportTicketMessage, error) {
	return s.supportTicketRepository.GetSupportTicketMessage(messageId)
}

func (s *SupportService) UpdateSupportTicketStatus(id int64, status int) error {
	return s.supportTicketRepository.UpdateSupportTicketStatus(id, status)
}

func (s *SupportService) GetSupportTickets() ([]model.SupportTicket, error) {
	return s.supportTicketRepository.GetSupportTickets()
}

func (s *SupportService) CloseTicket(id int64) error {
	return s.UpdateSupportTicketStatus(id, model.SupportTicketCloseStatus)
}

func (s *SupportService) getSupportTicketNotifierChannel(ticketId int64) string {
	return fmt.Sprintf("%s-%d", "[SUPPORT_TICKET]__", ticketId)
}

func (s *SupportService) NotifySupportTicketUpdate(ticketId int64) error {
	result := s.redisClient.Publish(context.TODO(), s.getSupportTicketNotifierChannel(ticketId), "1")
	return result.Err()
}

func (s *SupportService) ListenSupportTicketUpdate(ctx context.Context, ticketId int64, onNotified func() error) error {
	subscriber := s.redisClient.Subscribe(ctx, s.getSupportTicketNotifierChannel(ticketId))
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		if _, err := subscriber.ReceiveMessage(ctx); err != nil {
			return err
		}
		if err := onNotified(); err != nil {
			return err
		}
	}
}
