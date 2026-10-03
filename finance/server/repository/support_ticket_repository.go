package repository

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/server/model"
	"time"
)

type SupportTicketRepository struct {
	db *database.DB
}

func NewSupportTicketRepository(db *database.DB) *SupportTicketRepository {
	return &SupportTicketRepository{db}
}

func (r *SupportTicketRepository) GetSupportTickets() ([]model.SupportTicket, error) {
	var ret []model.SupportTicket
	if err := r.db.Select(&ret, "SELECT * FROM supportTickets WHERE updatedAt >= ? ORDER BY updatedAt DESC", time.Now().Unix()-(86400*30)); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *SupportTicketRepository) GetSupportTicketByID(id int64) (model.SupportTicket, error) {
	var ret model.SupportTicket
	if err := r.db.Get(&ret, "SELECT * FROM supportTickets WHERE id = ?", id); err != nil {
		return model.SupportTicket{}, err
	}
	return ret, nil
}

func (r *SupportTicketRepository) GetSupportTicketsByUserID(userId int64) ([]model.SupportTicket, error) {
	var ret []model.SupportTicket
	if err := r.db.Select(&ret, "SELECT * FROM supportTickets WHERE userId = ? ORDER BY updatedAt DESC", userId); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *SupportTicketRepository) GetSupportTicketMessages(ticketId int64) ([]model.SupportTicketMessage, error) {
	var ret []model.SupportTicketMessage
	if err := r.db.Select(&ret, "SELECT * FROM supportTicketMessages WHERE ticketId = ?", ticketId); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *SupportTicketRepository) GetSupportTicketMessagesAfterID(ticketId int64, afterId int64) ([]model.SupportTicketMessage, error) {
	var ret []model.SupportTicketMessage
	if err := r.db.Select(&ret, "SELECT * FROM supportTicketMessages WHERE ticketId = ? AND id > ?", ticketId, afterId); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *SupportTicketRepository) CreateSupportTicket(create model.CreateSupportTicket) (int64, error) {
	insert, err := r.db.Exec("INSERT INTO supportTickets (userId, category, createdAt, updatedAt) VALUES (?, ?, ?, ?)", create.UserID, create.Category, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *SupportTicketRepository) CreateSupportTicketMessage(create *model.CreateSupportTicketMessage) (int64, error) {
	insert, err := r.db.Exec("INSERT INTO supportTicketMessages (ticketId, author, role, message, attachments, createdAt) VALUES (?, ?, ?, ?, ?, ?)", create.TicketID, create.Author, create.Role, create.Message, create.RawAttachment_, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *SupportTicketRepository) UpdateSupportTicketStatus(id int64, status int) error {
	_, err := r.db.Exec("UPDATE supportTickets SET status = ?, updatedAt = ? WHERE id = ?", status, time.Now().Unix(), id)
	return err
}

func (r *SupportTicketRepository) GetSupportTicketMessage(messageId int64) (model.SupportTicketMessage, error) {
	var ret model.SupportTicketMessage
	if err := r.db.Get(&ret, "SELECT * FROM supportTicketMessages WHERE id = ?", messageId); err != nil {
		return model.SupportTicketMessage{}, err
	}
	return ret, nil
}
