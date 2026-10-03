package service

import (
	"fmt"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"time"
)

type PurchaseProcessorService struct {
	log             logrus.FieldLogger
	purchaseService *PurchaseService
	providers       map[string]provider.PurchaseProvider
}

func NewPurchaseProcessorService(log logrus.FieldLogger, scheduler gocron.Scheduler) *PurchaseProcessorService {
	service := &PurchaseProcessorService{
		log:       log,
		providers: make(map[string]provider.PurchaseProvider),
	}

	if _, err := scheduler.NewJob(gocron.DurationJob(time.Minute*10), gocron.NewTask(func() {
		if err := service.CheckStatusAll(); err != nil {
			log.Errorf("error periodic checking purchase status: %v", err)
		}
	}), gocron.WithName("purchaseStatusChecker"), gocron.WithStartAt(gocron.WithStartImmediately())); err != nil {
		log.Fatalf("error scheduling purchase status checker: %v", err)
	}

	return service
}

func (s *PurchaseProcessorService) SetPurchaseService(purchaseService *PurchaseService) {
	s.purchaseService = purchaseService
}

func (s *PurchaseProcessorService) RegisterProvider(skuPrefix string, provider provider.PurchaseProvider) {
	s.providers[skuPrefix] = provider
}

func (s *PurchaseProcessorService) ProcessPurchase(product model.Product, purchase model.Purchase) (*model.PurchaseProcessResult, error) {
	pr, ok := s.providers[product.Provider]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", product.Provider)
	}

	log := s.log.WithFields(logrus.Fields{
		"productId":   product.ID,
		"sku":         product.Sku,
		"provider":    product.Provider,
		"externalSku": product.ExternalSku,
		"purchaseId":  purchase.ID,
		"destination": purchase.Destination,
	})
	log.Infof("processing purchase...")

	res, err := pr.Process(product, purchase)
	if err != nil {
		log.Errorf("error processing purchase: %v", err)
		return nil, err
	}

	switch res.Status {
	case model.PurchaseStatusSuccess:
		log.Infof("purchase success")
	case model.PurchaseStatusFailed:
		log.Infof("purchase failed. reason: %v", res.FailedReason)
	case model.PurchaseStatusProcess:
		log.Infof("purchase is still in process")
	default:
		log.Errorf("invalid purchase status: %v", res.Status)
	}

	return res, nil
}

func (s *PurchaseProcessorService) CheckBill(product model.Product, destination string, refId string) (*model.PurchaseCheckBillResult, error) {
	pr, ok := s.providers[product.Provider]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", product.Provider)
	}

	return pr.CheckBill(product, destination, refId)
}

func (s *PurchaseProcessorService) HandleCallback(provider string, req *fasthttp.RequestCtx) error {
	pr, ok := s.providers[provider]
	if !ok {
		return fmt.Errorf("provider not found: %s", provider)
	}

	result, err := pr.HandleCallback(req)
	if err != nil {
		return err
	}

	if result == nil {
		return nil
	}

	purchase, err := s.purchaseService.GetPurchaseByRefId(result.RefId)
	if err != nil {
		return err
	}

	if result.Status == model.PurchaseStatusProcess {
		return nil
	}

	switch result.Status {
	case model.PurchaseStatusSuccess:
		return s.purchaseService.SetPurchaseSuccess(purchase, result.Proof)
	case model.PurchaseStatusFailed:
		return s.purchaseService.UpdatePurchaseStatus(purchase, &result.Status, "Failed", result.Proof)
	}

	return nil
}

func (s *PurchaseProcessorService) CheckStatus(purchase model.Purchase) error {
	if purchase.CreatedAt < (time.Now().Unix() - 2593000) { // 30 days
		return fmt.Errorf("purchase is too old, created at %s", time.Unix(purchase.CreatedAt, 0).Format(time.RFC822))
	}
	if purchase.Status != model.PurchaseStatusProcess {
		return fmt.Errorf("purchase is not in process")
	}

	pr, ok := s.providers[purchase.ProductProvider]
	if !ok {
		return fmt.Errorf("provider not found: %s", purchase.ProductProvider)
	}

	result, err := pr.CheckStatus(purchase)
	if err != nil {
		return err
	}

	if result.Status == model.PurchaseStatusProcess {
		return nil
	}

	switch result.Status {
	case model.PurchaseStatusSuccess:
		return s.purchaseService.SetPurchaseSuccess(purchase, result.Proof)
	case model.PurchaseStatusFailed:
		return s.purchaseService.UpdatePurchaseStatus(purchase, &result.Status, "Failed", result.Proof)
	}

	return nil
}

func (s *PurchaseProcessorService) CheckStatusAll() error {
	purchases, err := s.purchaseService.GetAllPurchasesByStatus(model.PurchaseStatusProcess)
	if err != nil {
		return err
	}

	for _, purchase := range purchases {
		// only allow purchases that >= 5 minutes and <= 3 days
		if purchase.CreatedAt < (time.Now().Unix()-300) || purchase.CreatedAt > (time.Now().Unix()-(259200)) {
			continue
		}
		if err := s.CheckStatus(purchase); err != nil {
			return err
		}
	}

	return nil
}
