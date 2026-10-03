package repository

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type PurchaseRepository struct {
	db *database.DB
}

func NewPurchaseRepository(db *database.DB) *PurchaseRepository {
	return &PurchaseRepository{
		db: db,
	}
}

func (r *PurchaseRepository) Begin() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *PurchaseRepository) GetPurchase(id int64) (model.Purchase, error) {
	var purchase model.Purchase
	if err := r.db.Get(&purchase, "SELECT * FROM purchases WHERE id = ?", id); err != nil {
		return model.Purchase{}, err
	}
	return purchase, nil
}

func (r *PurchaseRepository) GetPurchasesByUserId(userId int64, opt model.GetUserPurchasesOptions) ([]model.Purchase, error) {
	var purchases []model.Purchase
	query := "SELECT * FROM purchases WHERE userId = ?"
	params := []any{userId}
	if opt.Status != nil {
		query += " AND status = ?"
		params = append(params, *opt.Status)
	}
	if opt.BeforeID != 0 {
		query += " AND id < ?"
		params = append(params, opt.BeforeID)
	}
	query += " ORDER BY id DESC"
	if opt.Limit != 0 {
		query += " LIMIT ?"
		params = append(params, opt.Limit)
	}
	if err := r.db.Select(&purchases, query, params...); err != nil {
		return nil, err
	}
	return purchases, nil
}

func (r *PurchaseRepository) GetPurchaseByRefId(refId string) (model.Purchase, error) {
	var purchase model.Purchase
	if err := r.db.Get(&purchase, "SELECT * FROM purchases WHERE refId = ?", refId); err != nil {
		return model.Purchase{}, err
	}
	return purchase, nil
}

func (r *PurchaseRepository) UpdatePurchaseStatusWithTx(tx *sqlx.Tx, purchaseId int64, proof string, status int) error {
	_, err := tx.Exec("UPDATE purchases SET status = ?, proof = ? WHERE id = ?", status, proof, purchaseId)
	return err
}

func (r *PurchaseRepository) SetPurchaseSuccessWithTx(tx *sqlx.Tx, purchaseId int64, proof string, parsedProof string) error {
	_, err := tx.Exec("UPDATE purchases SET status = ?, proof = ?, parsedProof = ?, successAt = ? WHERE id = ?", model.PurchaseStatusSuccess, proof, parsedProof, time.Now().Unix(), purchaseId)
	return err
}

func (r *PurchaseRepository) SetPurchaseRefunded(tx *sqlx.Tx, purchaseId int64, refunded bool) error {
	_, err := tx.Exec("UPDATE purchases SET refunded = ? WHERE id = ?", refunded, purchaseId)
	return err
}

func (r *PurchaseRepository) AddPurchaseTrackWithTx(tx *sqlx.Tx, purchaseId int64, newStatus *int, description string) error {
	_, err := tx.Exec("INSERT INTO purchaseTracks (purchaseId, newStatus, description, createdAt) VALUES (?, ?, ?, ?)", purchaseId, newStatus, description, time.Now().Unix())
	return err
}

func (r *PurchaseRepository) CreatePurchaseWithTx(tx *sqlx.Tx, create model.CreatePurchase) (int64, error) {
	detailedDestination, err := create.DetailedDestination()
	if err != nil {
		return 0, err
	}
	readableDestination, err := create.ReadableDestination()
	if err != nil {
		return 0, err
	}
	extraData, err := create.ExtraData()
	if err != nil {
		return 0, err
	}
	billData, err := create.BillData()
	if err != nil {
		return 0, err
	}
	internalData, err := create.InternalData()
	if err != nil {
		return 0, err
	}
	insert, err := tx.Exec(
		"INSERT INTO purchases (userId, contactEmail, productId, productName, productKind, productType, productDescription, productSku, productCategoryId, productCategoryName, providerInfo, price, fee, wholesalePrice, destination, destinationId, detailedDestination, readableDestination, paymentMethod, status, refId, createdAt, paymentData, paymentRefId, paymentExpiredAt, extraData, internalData, totalBillFee, totalBill, billAdmin, billData, userSellPrice, note, parsedProof, productExternalSku, productProvider) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		create.UserID, create.ContactEmail, create.Product.ID, create.Product.Name, create.Product.Kind, create.Product.Type, create.Product.Description, create.Product.Sku, create.ProductCategory.ID, create.ProductCategory.Name, create.ProviderInfo, create.Price, create.Fee, create.WholesalePrice, create.Destination, create.Product.DestinationType, detailedDestination, readableDestination, create.PaymentMethod, model.PurchaseStatusWaitingPayment, create.RefId, time.Now().Unix(), create.PaymentData, create.PaymentRefId, create.PaymentExpiredAt, extraData, internalData, create.TotalBillFee, create.TotalBill, create.BillAdmin, billData, create.UserSellPrice, create.Product.PurchaseNote, "", create.Product.ExternalSku, create.Product.Provider)

	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *PurchaseRepository) GetLatestPurchases(limit int) ([]model.Purchase, error) {
	var ret []model.Purchase
	if err := r.db.Select(&ret, "SELECT * FROM purchases ORDER BY id DESC LIMIT ?", limit); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *PurchaseRepository) SearchPurchases(opt model.SearchPurchases) ([]model.Purchase, error) {
	var ret []model.Purchase
	if opt.DestinationQuery == "" {
		if err := r.db.Select(&ret, "SELECT * FROM purchases WHERE createdAt >= ? AND createdAt <= ? ORDER BY createdAt DESC", opt.FromDate, opt.ToDate); err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Select(&ret, "SELECT * FROM purchases WHERE createdAt >= ? AND createdAt <= ? AND destination LIKE ? ORDER BY createdAt DESC", opt.FromDate, opt.ToDate, "%"+opt.DestinationQuery+"%"); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (r *PurchaseRepository) GetPurchasesByDestinationLike(destination string) ([]model.Purchase, error) {
	var ret []model.Purchase
	if err := r.db.Select(&ret, "SELECT * FROM purchases WHERE destination LIKE ?", "%"+destination+"%"); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *PurchaseRepository) CreateBillPrePurchase(create model.CreateBillPrePurchase) (int64, error) {
	detailedDestination, err := create.DetailedDestination()
	if err != nil {
		return 0, err
	}
	readableDestination, err := create.ReadableDestination()
	if err != nil {
		return 0, err
	}
	extraData, err := create.ExtraData()
	if err != nil {
		return 0, err
	}
	billData, err := create.BillData()
	if err != nil {
		return 0, err
	}
	insert, err := r.db.Exec(
		"INSERT INTO billPrePurchases (userId, refId, contactEmail, productId, destination, destinationId, detailedDestination, readableDestination, billAmount, billAdmin, price, wholesalePrice, billData, extraData, createdAt, expiredAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		create.UserID, create.RefID, create.ContactEmail, create.ProductID, create.Destination, create.DestinationID, detailedDestination, readableDestination, create.BillAmount, create.BillAdmin, create.Price, create.WholesalePrice, billData, extraData, time.Now().Unix(), create.ExpiredAt)
	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *PurchaseRepository) GetBillPrePurchase(id int64) (model.BillPrePurchase, error) {
	var ret model.BillPrePurchase
	if err := r.db.Get(&ret, "SELECT * FROM billPrePurchases WHERE id = ?", id); err != nil {
		return model.BillPrePurchase{}, err
	}
	return ret, nil
}

func (r *PurchaseRepository) UpdatePurchaseUserSellPrice(id int64, userSellPrice int64) error {
	_, err := r.db.Exec("UPDATE purchases SET userSellPrice = ? WHERE id = ?", userSellPrice, id)
	return err
}

func (r *PurchaseRepository) GetAllPurchasesByStatus(status int) ([]model.Purchase, error) {
	var ret []model.Purchase
	if err := r.db.Select(&ret, "SELECT * FROM purchases WHERE status = ? ORDER BY id DESC", status); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *PurchaseRepository) GetUserPurchaseStats(userId int64, from int64, to int64) (model.UserPurchaseStats, error) {
	var ret model.UserPurchaseStats
	if err := r.db.Get(&ret, "SELECT COUNT(*) AS count, COALESCE(SUM(price), 0) AS total, COALESCE(SUM(IF(purchases.productType = ?, userSellPrice + billAdmin - totalBillFee, userSellPrice - price)), 0) AS profit FROM purchases WHERE userId = ? AND status = ? AND userSellPrice != 0 AND createdAt >= ? and createdAt <= ?", model.ProductTypeDigitalPostpaid, userId, model.PurchaseStatusSuccess, from, to); err != nil {
		return model.UserPurchaseStats{}, err
	}
	return ret, nil
}

func (r *PurchaseRepository) SearchPurchaseByUserId(id int64, query string) ([]model.Purchase, error) {
	var ret []model.Purchase
	if err := r.db.Select(&ret, "SELECT * FROM purchases WHERE userId = ? AND destination LIKE ? ORDER BY id DESC LIMIT 100", id, query+"%"); err != nil {
		return nil, err
	}
	return ret, nil
}
