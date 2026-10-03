package service

import (
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/repository"
)

type ProductService struct {
	productRepository *repository.ProductRepository
}

func NewProductService(productRepository *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (s *ProductService) GetProduct(productId int64) (model.Product, error) {
	return s.productRepository.GetProduct(productId)
}

func (s *ProductService) GetProductCategory(productCategoryId int64) (model.ProductCategory, error) {
	return s.productRepository.GetProductCategory(productCategoryId)
}

func (s *ProductService) GetAllProductCategoriesWithChildren(nested bool) ([]model.ProductCategory, error) {
	categories, err := s.productRepository.GetAllProductCategories()
	if err != nil {
		return nil, err
	}
	if !nested {
		return categories, nil
	}
	nestedCategories := make([]model.ProductCategory, 0)
	for _, cat := range categories {
		if cat.ParentID == nil {
			nestedCategories = append(nestedCategories, s.buildProductCategory(cat, categories))
		}
	}
	return nestedCategories, nil
}

func (s *ProductService) GetProductCategoryChildren(parentId int64, nested bool) ([]model.ProductCategory, error) {
	categories, err := s.productRepository.GetProductCategoryChildren(parentId)
	if err != nil {
		return nil, err
	}
	if !nested {
		return categories, nil
	}
	nestedCategories := make([]model.ProductCategory, 0)
	for _, cat := range categories {
		if *cat.ParentID == parentId {
			nestedCategories = append(nestedCategories, s.buildProductCategory(cat, categories))
		}
	}
	return nestedCategories, nil
}

func (s *ProductService) buildProductCategory(parent model.ProductCategory, categories []model.ProductCategory) model.ProductCategory {
	for _, cat := range categories {
		if cat.ParentID != nil && *cat.ParentID == parent.ID {
			cat = s.buildProductCategory(cat, categories)
			parent.Children = append(parent.Children, cat)
		}
	}
	return parent
}

func (s *ProductService) GetAllProducts() ([]model.Product, error) {
	return s.productRepository.GetAllProducts()
}

func (s *ProductService) GetProductByCategory(categoryId int64) ([]model.Product, error) {
	return s.productRepository.GetProductByCategory(categoryId)
}

func (s *ProductService) CreateProduct(create model.CreateOrUpdateProduct) error {
	return s.productRepository.CreateProduct(create)
}

func (s *ProductService) CreateCategory(create model.CreateOrUpdateProductCategory) error {
	return s.productRepository.CreateProductCategory(create)
}

func (s *ProductService) CreateProductDestination(create model.CreateOrUpdateProductDestination) error {
	return s.productRepository.CreateProductDestination(create)
}

func (s *ProductService) GetProductDestinations() ([]model.ProductDestination, error) {
	return s.productRepository.GetProductDestinations()
}

func (s *ProductService) GetProductDestination(id int64) (model.ProductDestination, error) {
	return s.productRepository.GetProductDestination(id)
}

func (s *ProductService) UpdateProductDestination(id int64, update model.CreateOrUpdateProductDestination) error {
	return s.productRepository.UpdateProductDestination(id, update)
}

func (s *ProductService) DeleteProductDestination(id int64) error {
	return s.productRepository.DeleteProductDestination(id)
}

func (s *ProductService) DeleteProduct(productId int64) error {
	return s.productRepository.DeleteProduct(productId)
}

func (s *ProductService) DeleteProductCategory(categoryId int64) error {
	return s.productRepository.DeleteProductCategory(categoryId)
}

func (s *ProductService) UpdateProduct(productId int64, update model.CreateOrUpdateProduct) error {
	return s.productRepository.UpdateProduct(productId, update)
}

func (s *ProductService) UpdateProductCategory(productId int64, update model.CreateOrUpdateProductCategory) error {
	return s.productRepository.UpdateProductCategory(productId, update)
}
