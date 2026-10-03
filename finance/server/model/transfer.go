package model

type TransferToUser struct {
	ID                    int64  `db:"id"`
	SenderUserID          int64  `db:"senderUserId"`
	SenderUsername        string `db:"senderUsername"`
	SenderName            string `db:"senderName"`
	SenderEmail           string `db:"senderEmail"`
	SenderTransactionID   int64  `db:"senderTransactionId"`
	ReceiverUserID        int64  `db:"receiverUserId"`
	ReceiverUsername      string `db:"receiverUsername"`
	ReceiverName          string `db:"receiverName"`
	ReceiverEmail         string `db:"receiverEmail"`
	ReceiverTransactionID int64  `db:"receiverTransactionId"`
	Purpose               int    `db:"purpose"`
	Amount                int64  `db:"amount"`
	Note                  string `db:"note"`
	CreatedAt             int64  `db:"createdAt"`
}

type CreateTransferToUser struct {
	SenderUser            User
	ReceiverUser          User
	SenderTransactionID   int64
	ReceiverTransactionID int64
	Purpose               int
	Amount                int64
	Note                  string
}
