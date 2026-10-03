package service

import (
	"github.com/akmalfairuz/finance/server/model"
	"time"
)

type AnalyticsService struct {
	authService     *AuthService
	paymentService  *PaymentService
	purchaseService *PurchaseService
	userService     *UserService
}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

func (s *AnalyticsService) SetPurchaseService(purchaseService *PurchaseService) {
	s.purchaseService = purchaseService
}

func (s *AnalyticsService) SetUserService(userService *UserService) {
	s.userService = userService
}

func (s *AnalyticsService) SetPaymentService(paymentService *PaymentService) {
	s.paymentService = paymentService
}

func (s *AnalyticsService) SetAuthService(authService *AuthService) {
	s.authService = authService
}

func (s *AnalyticsService) GenerateReport(from, to int64) (*model.AnalyticalReport, error) {
	ret := &model.AnalyticalReport{}

	purchases, err := s.purchaseService.SearchPurchases(model.SearchPurchases{
		FromDate: from,
		ToDate:   to,
	})
	if err != nil {
		return nil, err
	}
	for _, purchase := range purchases {
		switch purchase.Status {
		case model.PurchaseStatusProcess:
			ret.PendingPurchaseCount++
		case model.PurchaseStatusFailed:
			ret.FailedPurchaseCount++
		case model.PurchaseStatusSuccess:
			ret.TotalSuccessfulPurchasePrice += purchase.Price
			ret.TotalSuccessfulPurchaseWholesalePrice += purchase.WholesalePrice
			ret.SuccessfulPurchaseCount++
			ret.Profit += purchase.Price - purchase.WholesalePrice
		}
	}

	newUserCount, err := s.userService.GetNewUserCount(from, to)
	if err != nil {
		return nil, err
	}
	ret.NewUserCount = newUserCount

	totalUserBalance, err := s.userService.GetTotalUserBalance()
	if err != nil {
		return nil, err
	}
	ret.TotalUserBalance = totalUserBalance

	totalDeposit, err := s.paymentService.GetTotalDepositAmount(from, to)
	if err != nil {
		return nil, err
	}
	ret.TotalDeposit = totalDeposit

	depositCount, err := s.paymentService.GetAllDepositCount(from, to)
	if err != nil {
		return nil, err
	}
	ret.DepositCount = depositCount

	activeUserCount, err := s.authService.GetActiveUserCount(from, to)
	if err != nil {
		return nil, err
	}
	ret.ActiveUsers = activeUserCount

	ret.ReportDateStart = from
	ret.ReportDateEnd = to
	ret.CreatedAt = time.Now().Unix()

	return ret, nil
}
