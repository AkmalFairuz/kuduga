package responsetype

import "github.com/akmalfairuz/finance/server/model"

type UserPurchaseStatsWrapperResponse struct {
	PreviousMonth UserPurchaseStatsResponse `json:"previousMonth"`
	ThisMonth     UserPurchaseStatsResponse `json:"thisMonth"`
	Yesterday     UserPurchaseStatsResponse `json:"yesterday"`
	Today         UserPurchaseStatsResponse `json:"today"`
}

type UserPurchaseStatsResponse struct {
	Count  int64 `json:"count"`
	Profit int64 `json:"profit"`
	Total  int64 `json:"total"`
}

func NewUserPurchaseStatsResponse(p model.UserPurchaseStats) UserPurchaseStatsResponse {
	return UserPurchaseStatsResponse{
		Count:  p.Count,
		Profit: p.Profit,
		Total:  p.Total,
	}
}
