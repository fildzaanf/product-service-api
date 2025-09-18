package port // inbound

import (
	entity "product-service-api/internal/product/domain"
)

type ProductCommandServiceInterface interface {
	CreateProduct(product entity.Product, imageBytes []byte, imageFilename string) (entity.Product, error)
	UpdateProductByID(id string, product entity.Product, imageBytes []byte, imageFilename string) (entity.Product, error)
	DeleteProductByID(id string) error
}

type ProductQueryServiceInterface interface {
	GetProductByID(id string) (entity.Product, error)
	GetAllProducts() ([]entity.Product, error)
}
