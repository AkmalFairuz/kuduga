package service

import (
	"errors"
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/module/helper"
	"github.com/akmalfairuz/finance/module/pointer"
	"github.com/akmalfairuz/finance/module/random"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/akmalfairuz/finance/server/repository"
	"github.com/akmalfairuz/finance/server/service/proofparser"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type PurchaseService struct {
	log                      logrus.FieldLogger
	productService           *ProductService
	transactionService       *TransactionService
	paymentService           *PaymentService
	purchaseRepository       *repository.PurchaseRepository
	purchaseProcessorService *PurchaseProcessorService
	notificationService      *NotificationService
	checkerService           *CheckerService
	alert                    provider.TextAlertProvider
}

func NewPurchaseService(log logrus.FieldLogger, purchaseRepository *repository.PurchaseRepository) *PurchaseService {
	return &PurchaseService{
		log:                log,
		purchaseRepository: purchaseRepository,
	}
}

func (s *PurchaseService) SetPurchaseProcessorService(purchaseProcessorService *PurchaseProcessorService) {
	s.purchaseProcessorService = purchaseProcessorService
}

func (s *PurchaseService) SetProductService(productService *ProductService) {
	s.productService = productService
}

func (s *PurchaseService) SetPaymentService(paymentService *PaymentService) {
	s.paymentService = paymentService
}

func (s *PurchaseService) SetTransactionService(transactionService *TransactionService) {
	s.transactionService = transactionService
}

func (s *PurchaseService) SetNotificationService(notificationService *NotificationService) {
	s.notificationService = notificationService
}

func (s *PurchaseService) SetCheckerService(checkerService *CheckerService) {
	s.checkerService = checkerService
}

func (s *PurchaseService) SetAlert(alert provider.TextAlertProvider) {
	s.alert = alert
}

func (s *PurchaseService) generateRefId() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return random.StringWithChars(24, chars)
}

func (s *PurchaseService) GetBillPrePurchase(id int64) (model.BillPrePurchase, error) {
	return s.purchaseRepository.GetBillPrePurchase(id)
}

func (s *PurchaseService) CheckBill(userId *int64, email string, destination map[string]string, productId int64) (model.PurchaseCheckBillResult, error) {
	product, err := s.productService.GetProduct(productId)
	if err != nil {
		return model.PurchaseCheckBillResult{}, err
	}
	if product.Type != model.ProductTypeDigitalPostpaid {
		return model.PurchaseCheckBillResult{}, errors.New("not a postpaid product")
	}
	if !product.IsAvailable {
		return model.PurchaseCheckBillResult{}, errortype.ErrPurchaseProductNotAvailable
	}
	//if product.IsMaintenance() {
	//	return model.PurchaseCheckBillResult{}, errors.New("product is under maintenance")
	//}

	productDestination, err := s.productService.GetProductDestination(product.DestinationType)
	if err != nil {
		return model.PurchaseCheckBillResult{}, err
	}

	formattedDestination, err := productDestination.Format(destination)
	if err != nil {
		return model.PurchaseCheckBillResult{}, err
	}

	refId := product.Provider + "-POSTPAID-" + s.generateRefId()

	checkBillResult, err := s.purchaseProcessorService.CheckBill(product, formattedDestination, refId)
	if err != nil {
		return model.PurchaseCheckBillResult{}, err
	}

	id, err := s.purchaseRepository.CreateBillPrePurchase(model.CreateBillPrePurchase{
		RefID:                refId,
		UserID:               userId,
		ContactEmail:         email,
		ProductID:            productId,
		Destination:          formattedDestination,
		DestinationID:        product.DestinationType,
		DetailedDestination_: destination,
		ReadableDestination_: productDestination.Readable(destination),
		BillAmount:           checkBillResult.BillAmount,
		BillAdmin:            checkBillResult.RealAdmin,
		Price:                checkBillResult.CustomerPrice,
		WholesalePrice:       checkBillResult.Price,
		ExtraData_:           checkBillResult.AsExtraData(),
		BillData_:            checkBillResult.BillData,
		ExpiredAt:            time.Now().Add(time.Minute * 5).Unix(),
	})
	if err != nil {
		return model.PurchaseCheckBillResult{}, err
	}

	checkBillResult.ID = id
	return *checkBillResult, nil
}

func (s *PurchaseService) alertPurchase(purchase model.Purchase) {
	extraText := ""
	if purchase.Status == model.PurchaseStatusSuccess {
		extraText += fmt.Sprintf("Wkt Proses: %d dtk\n", purchase.SuccessAt-purchase.CreatedAt)
	}
	if err := s.alert.Alert(fmt.Sprintf("[Purchase #%d]\n%s - %s\n%s (%s)\n%s\n%s%s\n\n\n%s\n%s\n-------------------------------", purchase.ID, purchase.ProductSku, purchase.ProductName, helper.FormatRupiah(purchase.Price), helper.FormatRupiah(purchase.Price-purchase.WholesalePrice), purchase.Destination, extraText, strings.ToUpper(purchase.StatusText()), purchase.ContactEmail, helper.CurrentUTC7Date())); err != nil {
		s.log.Errorf("error sending depositAlert: %+v", err)
	}
}

func (s *PurchaseService) CreatePurchase(opt *model.CreatePurchaseOptions) (model.CreatePurchaseResult, error) {
	userId := opt.UserID
	email := opt.Email
	destination := opt.Destination
	productId := opt.ProductID
	bill := opt.Bill
	paymentMethod := opt.PaymentMethod

	guest := false
	if userId == nil {
		guest = true
	}

	product, err := s.productService.GetProduct(productId)
	if err != nil {
		return model.CreatePurchaseResult{}, err
	}
	if product.Category().RequireKyc() && !opt.UserKycVerified {
		return model.CreatePurchaseResult{}, errortype.ErrPurchaseKycRequired
	}

	// check maxPrice
	if opt.MaxPrice > 0 && product.Price > opt.MaxPrice {
		return model.CreatePurchaseResult{}, errortype.ErrPurchasePriceChange
	}

	isPostpaid := product.Type == model.ProductTypeDigitalPostpaid
	if isPostpaid {
		if bill == nil {
			return model.CreatePurchaseResult{}, errors.New("can't purchase postpaid product without bill")
		}
		if time.Now().Unix() >= bill.ExpiredAt {
			return model.CreatePurchaseResult{}, errortype.ErrPurchaseBillExpired
		}
		if *bill.UserID != *userId || bill.ContactEmail != email {
			return model.CreatePurchaseResult{}, errors.New("invalid bill")
		}
	}
	if !product.IsAvailable {
		return model.CreatePurchaseResult{}, errortype.ErrPurchaseProductNotAvailable
	}

	//if product.IsMaintenance() {
	//	return model.CreatePurchaseResult{}, errors.New("product is under maintenance")
	//}

	productDestination, err := s.productService.GetProductDestination(product.DestinationType)
	if err != nil {
		return model.CreatePurchaseResult{}, err
	}
	formattedDestination, err := productDestination.Format(destination)
	if err != nil {
		return model.CreatePurchaseResult{}, err
	}

	// check duplicate
	if !opt.AllowDuplicate {
		latestPurchases, err := s.purchaseRepository.GetPurchasesByUserId(*userId, model.GetUserPurchasesOptions{
			Status: pointer.Make(model.PurchaseStatusSuccess),
			Limit:  5,
		})
		if err != nil {
			return model.CreatePurchaseResult{}, err
		}
		for _, latestPurchase := range latestPurchases {
			if latestPurchase.UserID == opt.UserID && latestPurchase.ContactEmail == opt.Email && latestPurchase.ProductID == productId && latestPurchase.Destination == formattedDestination && latestPurchase.CreatedAt > time.Now().Unix()-180 {
				return model.CreatePurchaseResult{}, errortype.ErrPurchaseDuplicate
			}
		}
	}

	checkerResult := CheckerResult{}
	if !isPostpaid {
		if productDestination.CheckerID != "" {
			checkerResult, _ = s.checkerService.Check(productDestination.CheckerID, destination)
		}
	}

	if isPostpaid {
		if formattedDestination != bill.Destination {
			return model.CreatePurchaseResult{}, errors.New("mismatch destination with bill prepurchase")
		}
	}

	refId := ""
	if isPostpaid {
		refId = bill.RefID
	} else {
		refId = product.Provider + "-" + s.generateRefId()
	}

	userSellPrice := calculateAutoUserSellPrice(product.Price, 2000) // TODO: get from user setting
	paymentRefId := ""
	paymentData := ""
	fee := int64(0)
	price := int64(0)
	wholesalePrice := int64(0)
	totalBill := int64(0)
	totalBillFee := int64(0)
	billAdmin := int64(0)
	var billData []map[string]string
	if isPostpaid {
		userSellPrice = 1000
		price = bill.Price
		wholesalePrice = bill.WholesalePrice
		billAdmin = bill.BillAdmin
	} else {
		price = product.Price
		wholesalePrice = product.WholesalePrice
	}
	expiredAt := time.Now().Add(time.Minute * 5).Unix()
	var extraData map[string]any
	if isPostpaid {
		extraData, err = bill.ExtraData()
		if err != nil {
			return model.CreatePurchaseResult{}, err
		}
		billData2, err := bill.BillData()
		if err != nil {
			return model.CreatePurchaseResult{}, err
		}
		totalBill = bill.BillAmount
		// TODO: add bill.BillFee instead doing this
		totalBillFee = bill.Price - bill.BillAmount
		billData = billData2
	} else {
		extraData = map[string]any{}
		for _, checkerResultParts := range checkerResult {
			if checkerResultParts[1] == "" {
				continue
			}
			extraData[checkerResultParts[0]] = checkerResultParts[1]
		}
	}

	tx, err := s.purchaseRepository.Begin()
	if err != nil {
		_ = tx.Rollback()
		return model.CreatePurchaseResult{}, err
	}

	if guest {
		if paymentMethod == nil {
			return model.CreatePurchaseResult{}, errors.New("payment method is required")
		}
		customerName, _, _ := strings.Cut(opt.Email, "@")
		payment, err := s.paymentService.CreatePayment(customerName, "PURCHASE", *paymentMethod, product.Price)
		if err != nil {
			return model.CreatePurchaseResult{}, err
		}
		paymentRefId = payment.ExternalPaymentId
		paymentData = payment.ExternalPaymentData
		fee += payment.Fee
		price = payment.Amount
		expiredAt = payment.ExpiredAt
	} else {
		txnId, err := s.transactionService.CreateTransactionWithTx(tx, model.CreateTransaction{
			UserID:      *userId,
			IsAdd:       false,
			Type:        model.TransactionTypePurchase,
			ExtraData:   refId,
			Description: fmt.Sprintf("Purchase product %d/%s", product.ID, product.Sku),
			Amount:      price,
		})
		if err != nil {
			_ = tx.Rollback()
			return model.CreatePurchaseResult{}, err
		}
		paymentRefId = strconv.Itoa(int(txnId))
	}

	create := model.CreatePurchase{
		Product:              product,
		ProductCategory:      product.Category(),
		UserID:               userId,
		RefId:                refId,
		ContactEmail:         email,
		Price:                price,
		WholesalePrice:       wholesalePrice,
		Fee:                  fee,
		TotalBill:            totalBill,
		BillAdmin:            billAdmin,
		UserSellPrice:        userSellPrice,
		TotalBillFee:         totalBillFee,
		BillData_:            billData,
		Destination:          formattedDestination,
		DetailedDestination_: destination,
		ReadableDestination_: productDestination.Readable(destination),
		PaymentMethod:        paymentMethod,
		PaymentRefId:         paymentRefId,
		PaymentData:          paymentData,
		PaymentExpiredAt:     expiredAt,
		ExtraData_:           extraData,
	}

	purchaseId, err := s.purchaseRepository.CreatePurchaseWithTx(tx, create)
	if err != nil {
		_ = tx.Rollback()
		return model.CreatePurchaseResult{}, err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return model.CreatePurchaseResult{}, err
	}

	if !guest {
		if _, err := s.HandlePurchasePaymentReceived(purchaseId); err != nil {
			return model.CreatePurchaseResult{}, err
		}
	}

	s.log.WithFields(logrus.Fields{
		"email":              email,
		"purchaseId":         purchaseId,
		"productId":          product.ID,
		"productSku":         product.Sku,
		"productExternalSku": product.ExternalSku,
		"productProvider":    product.Provider,
		"price":              price,
		"wholesalePrice":     wholesalePrice,
		"destination":        formattedDestination,
	}).Info("purchase created")

	return model.CreatePurchaseResult{
		PurchaseID:       purchaseId,
		Price:            price + fee,
		RefId:            refId,
		PaymentExpiredAt: expiredAt,
	}, nil
}

func (s *PurchaseService) GetLatestPurchases(limit int) ([]model.Purchase, error) {
	return s.purchaseRepository.GetLatestPurchases(limit)
}

func (s *PurchaseService) SearchPurchases(opt model.SearchPurchases) ([]model.Purchase, error) {
	return s.purchaseRepository.SearchPurchases(opt)
}

func (s *PurchaseService) SearchPurchasesByUserId(userId int64, query string) ([]model.Purchase, error) {
	return s.purchaseRepository.SearchPurchaseByUserId(userId, query)
}

func (s *PurchaseService) GetPurchasesByDestinationLike(destination string) ([]model.Purchase, error) {
	return s.purchaseRepository.GetPurchasesByDestinationLike(destination)
}

func (s *PurchaseService) GetPurchase(purchaseId int64) (model.Purchase, error) {
	return s.purchaseRepository.GetPurchase(purchaseId)
}

func (s *PurchaseService) GetPurchasesByUserId(userId int64, opt model.GetUserPurchasesOptions) ([]model.Purchase, error) {
	return s.purchaseRepository.GetPurchasesByUserId(userId, opt)
}

func (s *PurchaseService) GetPurchaseByRefId(refId string) (model.Purchase, error) {
	return s.purchaseRepository.GetPurchaseByRefId(refId)
}

// HandlePurchasePaymentReceived is called when payment is received
func (s *PurchaseService) HandlePurchasePaymentReceived(purchaseId int64) (*model.PurchaseProcessResult, error) {
	purchase, err := s.purchaseRepository.GetPurchase(purchaseId)
	if err != nil {
		return nil, err
	}

	if purchase.Status != model.PurchaseStatusWaitingPayment {
		return nil, errors.New("purchase is not waiting for payment")
	}

	return s.ProcessPurchase(purchase)
}

func (s *PurchaseService) SetPurchaseSuccess(purchase model.Purchase, proof string) error {
	product, err := s.productService.GetProduct(purchase.ProductID)
	if err != nil {
		return fmt.Errorf("failed to get product while the purchase status will success: %w", err)
	}
	parsedProof := proofparser.Parse(product.ProofParser, proof)
	if parsedProof == nil {
		parsedProof = pointer.Make("")
	}

	return s.UpdatePurchaseStatus(purchase, pointer.Make(model.PurchaseStatusSuccess), "Purchase Success", proof, *parsedProof)
}

func (s *PurchaseService) UpdatePurchaseStatus(purchase model.Purchase, status *int, description string, proof ...string) error {
	tx, err := s.purchaseRepository.Begin()
	if err != nil {
		return err
	}

	proof2 := ""
	if len(proof) > 0 {
		proof2 = proof[0]
	}
	parsedProof := ""
	if len(proof) > 1 && proof[1] != "" {
		parsedProof = proof[1]
	}

	if err := database.Tx(tx, func() error {
		if status != nil {
			switch *status {
			case model.PurchaseStatusSuccess:
				if len(proof) == 0 {
					return errors.New("proof is required")
				}
				if err := s.purchaseRepository.SetPurchaseSuccessWithTx(tx, purchase.ID, proof[0], parsedProof); err != nil {
					return err
				}
			default:
				if err := s.purchaseRepository.UpdatePurchaseStatusWithTx(tx, purchase.ID, proof2, *status); err != nil {
					return err
				}
			}
		}
		if err := s.purchaseRepository.AddPurchaseTrackWithTx(tx, purchase.ID, status, description); err != nil {
			return err
		}
		// refund
		if !purchase.Refunded && purchase.PaymentMethod == nil && status != nil && *status == model.PurchaseStatusFailed {
			if _, err := s.transactionService.CreateTransactionWithTx(tx, model.CreateTransaction{
				UserID:      *purchase.UserID,
				Amount:      purchase.Price,
				IsAdd:       true,
				Description: fmt.Sprintf("Refund for purchase #%d", purchase.ID),
				ExtraData:   strconv.Itoa(int(purchase.ID)),
			}); err != nil {
				return err
			}
			if err := s.purchaseRepository.SetPurchaseRefunded(tx, purchase.ID, true); err != nil {
				return err
			}
		}
		if purchase.UserID != nil && status != nil && *status != model.PurchaseStatusProcess {
			_ = s.notificationService.PushNotificationToUser(*purchase.UserID, &provider.PushNotificationOptions{
				Notification: &provider.NotificationOptions{
					Title: fmt.Sprintf("Pembelian %s", model.PurchaseStatusText(*status)),
					Body:  fmt.Sprintf("Produk: %s - %s\nHarga: %s", purchase.ProductCategoryName, purchase.ProductName, helper.FormatRupiah(purchase.Price)),
				},
				Data: map[string]string{
					"_event": "purchaseStatusUpdate",
					"id":     strconv.Itoa(int(purchase.ID)),
				},
			})
		}
		return nil
	}); err != nil {
		return err
	}

	if status != nil {
		purchase.Status = *status
	}
	if purchase.Status == model.PurchaseStatusSuccess {
		purchase.ParsedProof = parsedProof
		purchase.Proof = &proof2
		purchase.SuccessAt = time.Now().Unix()
	}

	if s.alert != nil {
		go s.alertPurchase(purchase)
	}

	s.log.WithFields(logrus.Fields{
		"purchaseId":         purchase.ID,
		"productId":          purchase.ProductID,
		"productSku":         purchase.ProductSku,
		"productExternalSku": purchase.ProductExternalSku,
		"productProvider":    purchase.ProductProvider,
		"proof":              proof2,
		"price":              purchase.Price,
		"wholesalePrice":     purchase.WholesalePrice,
		"destination":        purchase.Destination,
	}).Infof("purchase status updated to %s", model.PurchaseStatusText(*status))

	return nil
}

func (s *PurchaseService) ProcessPurchase(purchase model.Purchase) (*model.PurchaseProcessResult, error) {
	product, err := s.productService.GetProduct(purchase.ProductID)
	if err != nil {
		return nil, err
	}

	if !product.IsAvailable {
		return nil, errortype.ErrPurchaseProductNotAvailable
	}

	//if product.IsMaintenance() {
	//	return nil, errors.New("product is under maintenance")
	//}

	if purchase.Status == model.PurchaseStatusWaitingPayment {
		if err := s.UpdatePurchaseStatus(purchase, pointer.Make(model.PurchaseStatusProcess), "Processing purchase"); err != nil {
			return nil, err
		}
	}

	result, err := s.purchaseProcessorService.ProcessPurchase(product, purchase)
	if err != nil {
		return nil, err
	}

	switch result.Status {
	case model.PurchaseStatusSuccess:
		if err := s.SetPurchaseSuccess(purchase, result.Proof); err != nil {
			return nil, err
		}
		return result, nil
	case model.PurchaseStatusFailed:
		if err := s.UpdatePurchaseStatus(purchase, pointer.Make(model.PurchaseStatusFailed), "Purchase failed", result.Proof); err != nil {
			return nil, err
		}
		return result, nil
	}

	return result, nil
}

func (s *PurchaseService) UpdatePurchaseUserSellPrice(purchase model.Purchase, userSellPrice int64) error {
	if err := s.purchaseRepository.UpdatePurchaseUserSellPrice(purchase.ID, userSellPrice); err != nil {
		return err
	}
	return nil
}

func (s *PurchaseService) GetAllPurchasesByStatus(status int) ([]model.Purchase, error) {
	return s.purchaseRepository.GetAllPurchasesByStatus(status)
}

func (s *PurchaseService) GetUserPurchaseStats(userId int64, from int64, to int64) (model.UserPurchaseStats, error) {
	return s.purchaseRepository.GetUserPurchaseStats(userId, from, to)
}

func calculateAutoUserSellPrice(price int64, targetAddPrice int64) int64 {
	if price%1000 >= 990 {
		price += 10
	}
	ret := price + targetAddPrice
	ret = ret - (ret % 1000)
	if price%1000 >= 500 {
		ret += 500
	}
	return ret
}
