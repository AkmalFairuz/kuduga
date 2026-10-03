package model

import "strings"

const (
	SupportTicketOpenStatus = iota
	SupportTicketCloseStatus
)

const (
	SupportTicketUserRole = iota
	SupportTicketCustomerServiceRole
	SupportTicketSystemRole
)

var SupportTicketCategories = []SupportTicketCategory{
	{ID: "purchase-invalid", Label: "Pembelian - status berhasil tapi tidak masuk"},
	{ID: "purchase-process-too-long", Label: "Pembelian - proses lama"},
	{ID: "purchase-other", Label: "Pembelian - lainnya"},
	{ID: "balance-deposit", Label: "Saldo - Deposit"},
	{ID: "balance-transfer", Label: "Saldo - Transfer"},
	{ID: "app", Label: "Aplikasi"},
	{ID: "account", Label: "Akun"},
	{ID: "other", Label: "Lainnya"},
}

type SupportTicketCategory struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type SupportTicket struct {
	ID         int64  `db:"id"`
	ExternalID string `db:"externalId"`
	UserID     int64  `db:"userId"`
	Category   string `db:"category"`
	Status     int    `db:"status"`
	CreatedAt  int64  `db:"createdAt"`
	UpdatedAt  int64  `db:"updatedAt"`
}

type SupportTicketMessage struct {
	ID           int64  `db:"id"`
	Author       string `db:"author"`
	TicketID     int64  `db:"ticketId"`
	Role         int    `db:"role"`
	Message      string `db:"message"`
	Attachments_ string `db:"attachments"`
	CreatedAt    int64  `db:"createdAt"`
}

func (model SupportTicketMessage) Attachments(prefix string) []string {
	if len(model.Attachments_) == 0 {
		return []string{}
	}
	parts := strings.Split(model.Attachments_, "\n")
	for i, part := range parts {
		parts[i] = prefix + "/" + part
	}
	return parts
}

type CreateSupportTicket struct {
	UserID   int64
	Category string
}

type CreateSupportTicketMessage struct {
	TicketID    int64
	Author      string
	Role        int
	Message     string
	Attachments []*File

	// RawAttachment_ do not define this
	RawAttachment_ string
}
