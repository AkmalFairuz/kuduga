package responsetype

import "github.com/akmalfairuz/finance/server/model"

type CreateSupportTicketResponse struct {
	TicketID int64 `json:"ticketId"`
}

type SupportTicketResponse struct {
	ID         int64                          `json:"id"`
	UserID     int64                          `json:"userId"`
	CategoryID string                         `json:"categoryId"`
	Category   string                         `json:"category"`
	Status     int                            `json:"status"`
	CreatedAt  int64                          `json:"createdAt"`
	UpdatedAt  int64                          `json:"updatedAt"`
	Messages   []SupportTicketMessageResponse `json:"messages,omitempty"`
}

func NewSupportTicketResponseFromModel(m model.SupportTicket) SupportTicketResponse {
	category := m.Category
	for _, cat := range model.SupportTicketCategories {
		if cat.ID == m.Category {
			category = cat.Label
			break
		}
	}
	return SupportTicketResponse{
		ID:         m.ID,
		UserID:     m.UserID,
		CategoryID: m.Category,
		Status:     m.Status,
		Category:   category,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

type SupportTicketMessageResponse struct {
	ID          int64    `json:"id"`
	Author      string   `json:"author"`
	Role        int      `json:"role"`
	Message     string   `json:"message"`
	Attachments []string `json:"attachments"`
	CreatedAt   int64    `json:"createdAt"`
}

func NewSupportTicketMessageResponseFromModel(m model.SupportTicketMessage, attachmentPrefix string) SupportTicketMessageResponse {
	return SupportTicketMessageResponse{
		ID:          m.ID,
		Author:      m.Author,
		Role:        m.Role,
		Message:     m.Message,
		CreatedAt:   m.CreatedAt,
		Attachments: m.Attachments(attachmentPrefix),
	}
}
