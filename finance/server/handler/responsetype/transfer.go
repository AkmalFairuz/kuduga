package responsetype

import "github.com/akmalfairuz/finance/server/model"

type BalanceTransferToUserResponse struct {
	TransferID int64 `json:"transferId"`
}

type BalanceTransferToUserHistoryResponse struct {
	ID        int64                           `json:"id"`
	IsSend    bool                            `json:"isSend"`
	Sender    BalanceTransferUserInfoResponse `json:"sender"`
	Receiver  BalanceTransferUserInfoResponse `json:"receiver"`
	Amount    int64                           `json:"amount"`
	Note      string                          `json:"note"`
	CreatedAt int64                           `json:"createdAt"`
}

func NewBalanceTransferToUserHistoryResponseFromModel(userId int64, model model.TransferToUser) BalanceTransferToUserHistoryResponse {
	return BalanceTransferToUserHistoryResponse{
		ID:     model.ID,
		IsSend: model.SenderUserID == userId,
		Sender: BalanceTransferUserInfoResponse{
			Username: model.SenderUsername,
			Name:     model.SenderName,
			Email:    model.SenderEmail,
		},
		Receiver: BalanceTransferUserInfoResponse{
			Username: model.ReceiverUsername,
			Name:     model.ReceiverName,
			Email:    model.ReceiverEmail,
		},
		Amount:    model.Amount,
		Note:      model.Note,
		CreatedAt: model.CreatedAt,
	}
}

type BalanceTransferUserInfoResponse struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}
