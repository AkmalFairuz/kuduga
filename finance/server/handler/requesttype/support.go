package requesttype

type CreateSupportTicketRequest struct {
	CategoryID string `form:"categoryId" validate:"required"`
	Message    string `form:"message" validate:"required,max=2000"`
}

type CreateSupportTicketMessageRequest struct {
	TicketID          int64  `form:"ticketId" validate:"required"`
	IsCustomerService bool   `form:"isCustomerService"`
	Message           string `form:"message" validate:"required,max=2000"`
}

type GetTicketDetailRequest struct {
	ID             int64 `query:"id" validate:"required"`
	MessageAfterID int64 `query:"messageAfterId"`
}
