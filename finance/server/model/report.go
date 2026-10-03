package model

type AnalyticalReport struct {
	ID int64 `db:"id"`

	Profit                                int64 `db:"profit"`
	NewUserCount                          int64 `db:"newUserCount"`
	ActiveUsers                           int64 `db:"activeUsers"`
	TotalDeposit                          int64 `db:"totalDeposit"`
	DepositCount                          int64 `db:"depositCount"`
	TotalUserBalance                      int64 `db:"totalUserBalance"`
	TotalSuccessfulPurchasePrice          int64 `db:"totalSuccessfulPurchasePrice"`
	TotalSuccessfulPurchaseWholesalePrice int64 `db:"totalSuccessfulPurchaseWholesalePrice"`
	SuccessfulPurchaseCount               int64 `db:"successfulPurchaseCount"`
	PendingPurchaseCount                  int64 `db:"pendingPurchaseCount"`
	FailedPurchaseCount                   int64 `db:"failedPurchaseCount"`

	ReportDateStart int64 `db:"reportDateStart"`
	ReportDateEnd   int64 `db:"reportDateEnd"`

	CreatedAt int64 `db:"createdAt"`
}
