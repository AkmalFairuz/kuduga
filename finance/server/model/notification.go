package model

import "encoding/json"

const (
	NotificationTypeUnknown = iota
	NotificationTypeSystem
	NotificationTypePurchase
	NotificationTypeDeposit
	NotificationTypeTransfer
	NotificationTypePromo
)

type Notification struct {
	ID          int64  `db:"id"`
	UserID      int64  `db:"userId"`
	Title       string `db:"title"`
	Description string `db:"description"`
	HasRead     bool   `db:"hasRead"`
	Type        int    `db:"type"`
	RefID       string `db:"refId"`
	Data_       string `db:"data"`
	CreatedAt   int64  `db:"createdAt"`
}

func (model Notification) Data() (map[string]string, error) {
	var ret map[string]string
	if err := json.Unmarshal([]byte(model.Data_), &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type CreateNotification struct {
	UserID      int64
	Title       string
	Description string
	Type        int
	RefID       string
	Data        map[string]string
	Push        bool
}

func (create CreateNotification) MarshalData() (string, error) {
	bytes, err := json.Marshal(create.Data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
