package repository

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/server/model"
	"time"
)

type NotificationRepository struct {
	db *database.DB
}

func NewNotificationRepository(db *database.DB) *NotificationRepository {
	return &NotificationRepository{db}
}

func (r *NotificationRepository) Create(create model.CreateNotification) (int64, error) {
	data, err := create.MarshalData()
	if err != nil {
		return 0, err
	}
	insert, err := r.db.Exec("INSERT INTO notifications (userId, hasRead, title, description, type, data, refId, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", create.UserID, false, create.Title, create.Description, create.Type, data, create.RefID, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *NotificationRepository) Get(userId int64) ([]model.Notification, error) {
	var ret []model.Notification
	if err := r.db.Select(&ret, "SELECT * FROM notifications WHERE userId = ? ORDER BY id DESC", userId); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *NotificationRepository) MarkRead(userId int64, id int64) error {
	_, err := r.db.Exec("UPDATE notifications SET hasRead = true WHERE id = ? AND userId = ?", userId, id)
	if err != nil {
		return err
	}
	return nil
}
