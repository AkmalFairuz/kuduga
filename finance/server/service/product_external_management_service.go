package service

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/digiflazz"
	"github.com/akmalfairuz/finance/module/helper"
	"github.com/akmalfairuz/finance/module/pointer"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
	"sort"
	"strings"
	"time"
)

type ProductExternalManagementService struct {
	digiflazzClient *digiflazz.Client
	productService  *ProductService
	log             logrus.FieldLogger
	alert           provider.TextAlertProvider
}

func NewProductExternalManagementService(log logrus.FieldLogger, scheduler gocron.Scheduler, digiflazzClient *digiflazz.Client) *ProductExternalManagementService {
	service := &ProductExternalManagementService{digiflazzClient: digiflazzClient, log: log}
	if _, err := scheduler.NewJob(gocron.DurationJob(time.Second*30), gocron.NewTask(func() {
		if err := service.SyncDigiflazzProducts(); err != nil {
			log.Errorf("error syncing digiflazz products: %+v", err)
		}
	}), gocron.WithName("digiflazzProductsSyncer")); err != nil {
		log.Fatalf("error creating digiflazzProductsSyncer scheduler: %+v", err)
	}
	return service
}

func (s *ProductExternalManagementService) SetProductService(productService *ProductService) {
	s.productService = productService
}

func (s *ProductExternalManagementService) SetAlert(alert provider.TextAlertProvider) {
	s.alert = alert
}

const (
	minimumCommission = 5
)

func (s *ProductExternalManagementService) SyncDigiflazzProducts() error {
	digiflazzProducts, err := s.digiflazzClient.GetProducts()
	if err != nil {
		return err
	}
	products, err := s.productService.GetAllProducts()
	if err != nil {
		return err
	}
	alerts := make([]string, 0)
	for _, product := range products {
		if product.Provider != "digiflazz" {
			continue
		}
		log := s.log.WithField("sku", product.Sku)
		found := false
		for _, dgp := range digiflazzProducts {
			if product.ExternalSku != dgp.GetBuyerSkuCode() {
				continue
			}

			update := model.CreateOrUpdateProduct{}
			if !product.IsAvailable && dgp.GetBuyerProductStatus() && dgp.GetSellerProductStatus() && product.MaxWholesalePrice >= dgp.GetPrice() {
				update.IsAvailable = pointer.Make(true)
			} else if product.IsAvailable && (!dgp.GetBuyerProductStatus() || !dgp.GetSellerProductStatus()) {
				break
			}
			if product.WholesalePrice != dgp.GetPrice() {
				if dgp.GetPrice() > product.MaxWholesalePrice {
					break
				}
				update.WholesalePrice = dgp.GetPrice()
				if product.Price-dgp.GetPrice() < minimumCommission { // minimum commission
					update.Price = dgp.GetPrice() + minimumCommission
				}
				alerts = append(alerts, fmt.Sprintf("[Wholesale Price Change #%d]\n[%s] %s - %s\n%s -> %s", product.ID, product.Sku, product.CategoryName, product.Name, helper.FormatRupiah(product.WholesalePrice), helper.FormatRupiah(update.WholesalePrice)))
				log.Infof("Wholesale Price updating from %d to %d. Price updating from %d to %d", product.WholesalePrice, update.WholesalePrice, product.Price, update.Price)
			}
			{
				dgp2, ok := dgp.(digiflazz.ProductResponse)
				if ok {
					// fix cut off
					hh, mm, _ := strings.Cut(dgp2.StartCutOff, ":")
					hh2, mm2, _ := strings.Cut(dgp2.EndCutOff, ":")
					hh = fmt.Sprintf("%02s", hh)
					hh2 = fmt.Sprintf("%02s", hh2)
					mm = fmt.Sprintf("%02s", mm)
					mm2 = fmt.Sprintf("%02s", mm2)
					dgp2.StartCutOff = fmt.Sprintf("%s:%s", hh, mm)
					dgp2.EndCutOff = fmt.Sprintf("%s:%s", hh2, mm2)

					if dgp2.StartCutOff == "00:00" && dgp2.EndCutOff == "00:00" && (product.CutOffStart != "" || product.CutOffEnd != "") {
						update.CutOffStart = pointer.Make("")
						update.CutOffEnd = pointer.Make("")
					} else {
						if product.CutOffStart != dgp2.StartCutOff {
							update.CutOffStart = &dgp2.StartCutOff
						}
						if product.CutOffEnd != dgp2.EndCutOff {
							update.CutOffEnd = &dgp2.EndCutOff
						}
					}
				}
			}

			{
				dgp2, ok := dgp.(digiflazz.ProductPostpaidResponse)
				if ok {
					if dgp2.Admin != product.BillAdmin {
						update.BillAdmin = &dgp2.Admin
					}
				}
			}

			if err := s.productService.UpdateProduct(product.ID, update); err != nil {
				return err
			}
			found = true
			break
		}
		if !found && product.IsAvailable {
			update := model.CreateOrUpdateProduct{}
			update.IsAvailable = pointer.Make(false)
			alerts = append(alerts, fmt.Sprintf("[Product Not Available #%d]\n[%s] %s - %s", product.ID, product.Sku, product.CategoryName, product.Name))
			log.Warnf("product is not available")
			if err := s.productService.UpdateProduct(product.ID, update); err != nil {
				return err
			}
		}
	}

	if s.alert != nil && len(alerts) > 0 {
		go func() {
			if err := s.alert.Alert(strings.Join(alerts, "\n<--------------------------->\n")); err != nil {
				s.log.Errorf("error sending depositAlert: %+v", err)
			}
		}()
	}
	return nil
}

func (s *ProductExternalManagementService) ImportFromDigiflazz(toCategoryId int64, productKind string, destinationId int64, addPrice int64, prefix, brand, productType, category string) (int64, error) {
	products, err := s.digiflazzClient.GetPrepaidPriceList("")
	if err != nil {
		return 0, err
	}

	sort.Slice(products, func(i, j int) bool {
		return products[i].Price+addPrice < products[j].Price+addPrice
	})

	existingProducts, err := s.productService.GetAllProducts()
	if err != nil {
		return 0, err
	}

	importedProducts := int64(0)
	for _, product := range products {
		alreadyExists := false
		for _, existingProduct := range existingProducts {
			if existingProduct.ExternalSku == product.BuyerSkuCode {
				alreadyExists = true
				break
			}
		}
		if alreadyExists {
			continue
		}
		if brand != "" && product.Brand != brand {
			continue
		}
		if productType != "" && product.Type != productType {
			continue
		}
		if category != "" && product.Category != category {
			continue
		}
		if !containsWildcard(product.ProductName, prefix) {
			continue
		}
		if product.StartCutOff == "00:00" && product.EndCutOff == "00:00" {
			product.StartCutOff = ""
			product.EndCutOff = ""
		}
		if err := s.productService.CreateProduct(model.CreateOrUpdateProduct{
			CategoryID:          toCategoryId,
			Sku:                 product.BuyerSkuCode,
			ExternalSku:         product.BuyerSkuCode,
			Provider:            provider.DigiflazzPurchaseProviderName,
			Name:                product.ProductName,
			Kind:                &productKind,
			SkuPriority:         pointer.Make(0),
			PurchaseNote:        pointer.Make(""),
			Description:         pointer.Make(""),
			Type:                model.ProductTypeDigital,
			CutOffStart:         &product.StartCutOff,
			CutOffEnd:           &product.EndCutOff,
			ImageUrl:            pointer.Make(""),
			DestinationType:     destinationId,
			IsAvailable:         pointer.Make(true),
			Price:               product.Price + addPrice,
			BillAdmin:           pointer.Make(int64(0)),
			MaxWholesalePrice:   product.Price + addPrice,
			WholesalePrice:      product.Price,
			BeforeDiscountPrice: pointer.Make(int64(0)),
			ProofParser:         pointer.Make(""),
		}); err != nil {
			return importedProducts, fmt.Errorf("error when importing product %s: %w", product.BuyerSkuCode, err)
		}
		importedProducts++
	}
	return importedProducts, nil
}

func containsWildcard(s, substr string) bool {
	if len(substr) == 0 {
		return true // If the substring is empty, consider it as a match
	}

	wildcardIndex := strings.Index(substr, "*")
	if wildcardIndex == -1 {
		return strings.Contains(s, substr) // If no wildcard, use strings.Contains directly
	}

	prefix := substr[:wildcardIndex]
	suffix := substr[wildcardIndex+1:]

	if len(prefix) > 0 && !strings.HasPrefix(s, prefix) {
		return false
	}

	if len(suffix) > 0 && !strings.HasSuffix(s, suffix) {
		return false
	}

	return true
}
