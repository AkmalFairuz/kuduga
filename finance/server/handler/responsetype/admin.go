package responsetype

import "github.com/akmalfairuz/finance/server/model"

type AdminProductResponse struct {
	ID                  int64                      `json:"id"`
	CategoryId          int64                      `json:"categoryId"`
	CategoryName        string                     `json:"categoryName"`
	Sku                 string                     `json:"sku"`
	SkuPriority         int                        `json:"skuPriority"`
	ExternalSku         string                     `json:"externalSku"`
	Provider            string                     `json:"provider"`
	Type                int                        `json:"type"`
	IsAvailable         bool                       `json:"isAvailable"`
	Status              string                     `json:"status"`
	ImageUrl            string                     `json:"imageUrl"`
	Name                string                     `json:"name"`
	Kind                string                     `json:"kind"`
	PurchaseNote        string                     `json:"purchaseNote"`
	Description         string                     `json:"description"`
	Destination         ProductDestinationResponse `json:"destination"`
	Price               int64                      `json:"price"`
	BillAdmin           int64                      `json:"billAdmin"`
	MaxWholesalePrice   int64                      `json:"maxWholesalePrice"`
	WholesalePrice      int64                      `json:"wholesalePrice"`
	BeforeDiscountPrice int64                      `json:"beforeDiscountPrice"`
	CutOffStart         string                     `json:"cutOffStart"`
	CutOffEnd           string                     `json:"cutOffEnd"`
	ProofParser         string                     `json:"proofParser"`
}

type AdminSummaryResponse struct {
	CurrentMonthProfit int64 `json:"thisMonthProfit"`
	TodayProfit        int64 `json:"todayProfit"`

	TodaySuccessfulPurchases int64 `json:"todaySuccessfulPurchases"`
	TodayFailedPurchases     int64 `json:"todayFailedPurchases"`

	TotalUsersBalance int64 `json:"totalUsersBalance"`

	ProviderBalance map[string]int64 `json:"providerBalance"`
}

type AdminPurchaseResponse struct {
	ID int64 `json:"id"`

	UserID       *int64 `json:"userId"`
	ContactEmail string `json:"email"`

	ProductId           int64  `json:"productId"`
	ProductSku          string `json:"productSku"`
	ProductName         string `json:"productName"`
	ProductDescription  string `json:"productDescription"`
	ProductCategoryName string `json:"productCategoryName"`

	Status     int    `json:"status"`
	StatusText string `json:"statusText"`

	Destination         string            `json:"destination"`
	DestinationId       int64             `json:"destinationId"`
	DetailedDestination map[string]string `json:"detailedDestination"`

	Proof *string `json:"proof"`

	WholesalePrice int64 `json:"wholesalePrice"`
	Price          int64 `json:"price"`

	Fee              int64  `json:"fee"`
	PaymentMethod    *int   `json:"paymentMethod"`
	PaymentRefId     string `json:"paymentRefId"`
	PaymentExpiredAt int64  `json:"paymentExpiredAt"`
	PaymentData      string `json:"paymentData"`

	SuccessAt int64 `json:"successAt"`
	CreatedAt int64 `json:"createdAt"`
}

func NewAdminPurchaseResponseFromModel(purchase model.Purchase) (AdminPurchaseResponse, error) {
	detailedDestination, err := purchase.DetailedDestination()
	if err != nil {
		return AdminPurchaseResponse{}, err
	}
	ret := AdminPurchaseResponse{
		ID:                  purchase.ID,
		UserID:              purchase.UserID,
		ContactEmail:        purchase.ContactEmail,
		ProductId:           purchase.ProductID,
		ProductSku:          purchase.ProductSku,
		ProductName:         purchase.ProductName,
		ProductDescription:  purchase.ProductDescription,
		ProductCategoryName: purchase.ProductCategoryName,
		Status:              purchase.Status,
		StatusText:          purchase.StatusText(),
		Destination:         purchase.Destination,
		DestinationId:       purchase.DestinationId,
		DetailedDestination: detailedDestination,
		Proof:               purchase.Proof,
		WholesalePrice:      purchase.WholesalePrice,
		Price:               purchase.Price,
		Fee:                 purchase.Fee,
		PaymentMethod:       purchase.PaymentMethod,
		PaymentRefId:        purchase.PaymentRefId,
		PaymentExpiredAt:    purchase.PaymentExpiredAt,
		PaymentData:         purchase.PaymentData,
		SuccessAt:           purchase.SuccessAt,
		CreatedAt:           purchase.CreatedAt,
	}
	return ret, nil
}

type AdminUserKycResponse struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	FullName    string `json:"fullName"`
	Status      int    `json:"status"`
	DocumentID  string `json:"documentId"`
	DocumentUrl string `json:"documentUrl"`
	CreatedAt   int64  `json:"createdAt"`
}

func NewAdminUserKycResponseFromModel(kyc model.UserKyc, attachmentPrefix string) AdminUserKycResponse {
	return AdminUserKycResponse{
		ID:          kyc.ID,
		UserID:      kyc.UserID,
		FullName:    kyc.FullName,
		Status:      kyc.Status,
		DocumentID:  kyc.DocumentID,
		DocumentUrl: attachmentPrefix + "/" + kyc.DocumentFileID,
		CreatedAt:   kyc.CreatedAt,
	}
}

type AdminMetaResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AdminBusinessReportResponse struct {
	ID int64 `json:"id"`

	Profit                                int64 `json:"profit"`
	NewUserCount                          int64 `json:"newUserCount"`
	ActiveUsers                           int64 `json:"activeUsers"`
	TotalDeposit                          int64 `json:"totalDeposit"`
	DepositCount                          int64 `json:"depositCount"`
	TotalUserBalance                      int64 `json:"totalUserBalance"`
	TotalSuccessfulPurchasePrice          int64 `json:"totalSuccessfulPurchasePrice"`
	TotalSuccessfulPurchaseWholesalePrice int64 `json:"totalSuccessfulPurchaseWholesalePrice"`
	SuccessfulPurchaseCount               int64 `json:"successfulPurchaseCount"`
	PendingPurchaseCount                  int64 `json:"pendingPurchaseCount"`
	FailedPurchaseCount                   int64 `json:"failedPurchaseCount"`

	ReportDateStart int64 `json:"reportDateStart"`
	ReportDateEnd   int64 `json:"reportDateEnd"`

	CreatedAt int64 `json:"createdAt"`
}

func NewAdminBusinessReportResponseFromModel(report model.AnalyticalReport) AdminBusinessReportResponse {
	return AdminBusinessReportResponse{
		ID:                                    report.ID,
		Profit:                                report.Profit,
		NewUserCount:                          report.NewUserCount,
		ActiveUsers:                           report.ActiveUsers,
		TotalDeposit:                          report.TotalDeposit,
		DepositCount:                          report.DepositCount,
		TotalUserBalance:                      report.TotalUserBalance,
		TotalSuccessfulPurchasePrice:          report.TotalSuccessfulPurchasePrice,
		TotalSuccessfulPurchaseWholesalePrice: report.TotalSuccessfulPurchaseWholesalePrice,
		SuccessfulPurchaseCount:               report.SuccessfulPurchaseCount,
		PendingPurchaseCount:                  report.PendingPurchaseCount,
		FailedPurchaseCount:                   report.FailedPurchaseCount,
		ReportDateStart:                       report.ReportDateStart,
		ReportDateEnd:                         report.ReportDateEnd,
		CreatedAt:                             report.CreatedAt,
	}
}
