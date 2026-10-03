package repository

import (
	"encoding/json"
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/server/model"
	"time"
)

type ProductRepository struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

var selectProductQuery = "SELECT products.*, productDestinations.name AS destinationName, productDestinations.description AS destinationDescription, productCategories.name AS categoryName, productCategories.parentId AS categoryParentId, productCategories.meta AS categoryMeta FROM products LEFT JOIN productCategories ON products.categoryId = productCategories.id LEFT JOIN productDestinations ON products.destinationType = productDestinations.id"
var selectProductByIdQuery = selectProductQuery + " WHERE products.id = ?"
var selectProductByCategoryIdQuery = selectProductQuery + " WHERE products.categoryId = ?"
var selectProductCategoriesChildrenQuery = `
WITH RECURSIVE SubCategoryHierarchy AS (
    SELECT id, parentId, name, description, iconUrl, productType, meta
    FROM productCategories
    WHERE parentId = ?
    UNION ALL
    SELECT pc.id, pc.parentId, pc.name, pc.description, pc.iconUrl, pc.productType, pc.meta
    FROM productCategories pc
    JOIN SubCategoryHierarchy sch ON pc.parentId = sch.id
)
SELECT id, parentId, name, description, iconUrl, productType, meta
FROM SubCategoryHierarchy
ORDER BY id;
`

func (r *ProductRepository) CreateProductCategory(create model.CreateOrUpdateProductCategory) error {
	if _, err := r.db.Exec("INSERT INTO productCategories (name, parentId, description, iconUrl, productType, meta) VALUES (?, ?, ?, ?, ?, ?)", create.Name, create.ParentID, create.Description, create.IconUrl, create.ProductType, create.Meta); err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetProductCategory(productCategoryId int64) (model.ProductCategory, error) {
	var product model.ProductCategory
	if err := r.db.Get(&product, "SELECT * FROM productCategories WHERE id = ?", productCategoryId); err != nil {
		return model.ProductCategory{}, err
	}
	return product, nil
}

func (r *ProductRepository) GetAllProductCategories() ([]model.ProductCategory, error) {
	var categories []model.ProductCategory
	if err := r.db.Select(&categories, "SELECT * FROM productCategories"); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductRepository) GetProductCategoryChildren(parentId int64) ([]model.ProductCategory, error) {
	var categories []model.ProductCategory
	if err := r.db.Select(&categories, selectProductCategoriesChildrenQuery, parentId); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductRepository) CreateProduct(create model.CreateOrUpdateProduct) error {
	now := time.Now().Unix()
	if _, err := r.db.Exec("INSERT INTO products (categoryId, sku, skuPriority, externalSku, provider, name, kind, type, imageUrl, cutOffStart, cutOffEnd, destinationType, description, purchaseNote, isAvailable, price, billAdmin, maxWholesalePrice, wholesalePrice, beforeDiscountPrice, proofParser, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		create.CategoryID, create.Sku, create.SkuPriority, create.ExternalSku, create.Provider, create.Name, *create.Kind, create.Type, *create.ImageUrl, *create.CutOffStart, *create.CutOffEnd, create.DestinationType, *create.Description, *create.PurchaseNote, *create.IsAvailable, create.Price, *create.BillAdmin, create.MaxWholesalePrice, create.WholesalePrice, *create.BeforeDiscountPrice, *create.ProofParser, now, now); err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetProductByCategory(categoryId int64) ([]model.Product, error) {
	var products []model.Product
	if err := r.db.Select(&products, selectProductByCategoryIdQuery, categoryId); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetProduct(productId int64) (model.Product, error) {
	var product model.Product
	if err := r.db.Get(&product, selectProductByIdQuery, productId); err != nil {
		return model.Product{}, err
	}
	return product, nil
}

func (r *ProductRepository) UpdateProduct(productId int64, update model.CreateOrUpdateProduct) error {
	now := time.Now().Unix()
	product, err := r.GetProduct(productId)
	if err != nil {
		return err
	}

	if update.CategoryID == 0 {
		update.CategoryID = product.CategoryID
	}
	if update.Sku == "" {
		update.Sku = product.Sku
	}
	if update.ExternalSku == "" {
		update.ExternalSku = product.ExternalSku
	}
	if update.Provider == "" {
		update.Provider = product.Provider
	}
	if update.Name == "" {
		update.Name = product.Name
	}
	if update.Description == nil {
		update.Description = &product.Description
	}
	if update.Type == 0 {
		update.Type = product.Type
	}
	if update.ImageUrl == nil {
		update.ImageUrl = &product.ImageUrl
	}
	if update.CutOffStart == nil {
		update.CutOffStart = &product.CutOffStart
	}
	if update.CutOffEnd == nil {
		update.CutOffEnd = &product.CutOffEnd
	}
	if update.IsAvailable == nil {
		update.IsAvailable = &product.IsAvailable
	}
	if update.Price == 0 {
		update.Price = product.Price
	}
	if update.BillAdmin == nil {
		update.BillAdmin = &product.BillAdmin
	}
	if update.MaxWholesalePrice == 0 {
		update.MaxWholesalePrice = product.MaxWholesalePrice
	}
	if update.WholesalePrice == 0 {
		update.WholesalePrice = product.WholesalePrice
	}
	if update.BeforeDiscountPrice == nil {
		update.BeforeDiscountPrice = &product.BeforeDiscountPrice
	}
	if update.Kind == nil {
		update.Kind = &product.Kind
	}
	if update.PurchaseNote == nil {
		update.PurchaseNote = &product.PurchaseNote
	}
	if update.DestinationType == 0 {
		update.DestinationType = product.DestinationType
	}
	if update.ProofParser == nil {
		update.ProofParser = &product.ProofParser
	}
	if update.SkuPriority == nil {
		update.SkuPriority = &product.SkuPriority
	}

	if _, err := r.db.Exec("UPDATE products SET categoryId = ?, externalSku = ?, sku = ?, skuPriority = ?, name = ?, kind = ?, type = ?, imageUrl = ?, cutOffStart = ?, cutOffEnd = ?, destinationType = ?, description = ?, purchaseNote = ?, isAvailable = ?, price = ?, billAdmin = ?, maxWholesalePrice = ?, wholesalePrice = ?, beforeDiscountPrice = ?, updatedAt = ?, proofParser = ? WHERE id = ?",
		update.CategoryID, update.ExternalSku, update.Sku, *update.SkuPriority, update.Name, *update.Kind, update.Type, *update.ImageUrl, *update.CutOffStart, *update.CutOffEnd, update.DestinationType, *update.Description, *update.PurchaseNote, update.IsAvailable, update.Price, *update.BillAdmin, update.MaxWholesalePrice, update.WholesalePrice, *update.BeforeDiscountPrice, now, *update.ProofParser, productId); err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) UpdateProductCategory(productCategoryId int64, update model.CreateOrUpdateProductCategory) error {
	productCategory, err := r.GetProductCategory(productCategoryId)
	if err != nil {
		return err
	}

	if update.Name == "" {
		update.Name = productCategory.Name
	}
	if update.ParentID == nil {
		update.ParentID = &productCategory.ParentID
	}
	if update.Description == nil {
		update.Description = &productCategory.Description
	}
	if update.IconUrl == nil {
		update.IconUrl = &productCategory.IconUrl
	}
	if update.ProductType == 0 {
		update.ProductType = productCategory.ProductType
	}
	if update.Meta == nil {
		update.Meta = &productCategory.Meta_
	}

	if _, err := r.db.Exec("UPDATE productCategories SET name = ?, parentId = ?, description = ?, iconUrl = ?, productType = ?, meta = ? WHERE id = ?", update.Name, *update.ParentID, *update.Description, *update.IconUrl, update.ProductType, *update.Meta, productCategoryId); err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(productId int64) error {
	if _, err := r.db.Exec("DELETE FROM products WHERE id = ?", productId); err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) DeleteProductCategory(categoryId int64) error {
	if _, err := r.db.Exec("DELETE FROM productCategories WHERE id = ?", categoryId); err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetAllProducts() ([]model.Product, error) {
	var products []model.Product
	if err := r.db.Select(&products, selectProductQuery); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) CreateProductDestination(create model.CreateOrUpdateProductDestination) error {
	fieldBytes, err := json.Marshal(create.Fields)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("INSERT INTO productDestinations (name, description, format, fields, checkerId) VALUES (?, ?, ?, ?, ?)", create.Name, create.Description, create.Format, string(fieldBytes), *create.CheckerID)
	return err
}

func (r *ProductRepository) GetProductDestinations() ([]model.ProductDestination, error) {
	var productDestinations []model.ProductDestination
	if err := r.db.Select(&productDestinations, "SELECT * FROM productDestinations"); err != nil {
		return nil, err
	}
	return productDestinations, nil
}

func (r *ProductRepository) GetProductDestination(id int64) (model.ProductDestination, error) {
	var productDestination model.ProductDestination
	if err := r.db.Get(&productDestination, "SELECT * FROM productDestinations WHERE id = ?", id); err != nil {
		return model.ProductDestination{}, err
	}
	return productDestination, nil
}

func (r *ProductRepository) UpdateProductDestination(id int64, update model.CreateOrUpdateProductDestination) error {
	destination, err := r.GetProductDestination(id)
	if err != nil {
		return err
	}
	if update.Name == "" {
		update.Name = destination.Name
	}
	if update.Description == nil {
		update.Description = &destination.Description
	}
	if update.Format == nil {
		update.Format = &destination.Format_
	}
	if update.CheckerID == nil {
		update.CheckerID = &destination.CheckerID
	}
	var fields string
	if update.Fields != nil {
		var fieldBytes []byte
		if fieldBytes, err = json.Marshal(update.Fields); err != nil {
			return err
		}
		fields = string(fieldBytes)
	} else {
		fields = destination.Fields_
	}
	_, err = r.db.Exec("UPDATE productDestinations SET name = ?, description = ?, format = ?, fields = ?, checkerId = ? WHERE id = ?", update.Name, update.Description, update.Format, fields, update.CheckerID, id)
	return err
}

func (r *ProductRepository) DeleteProductDestination(id int64) error {
	_, err := r.db.Exec("DELETE FROM productDestinations WHERE id = ?", id)
	return err
}
