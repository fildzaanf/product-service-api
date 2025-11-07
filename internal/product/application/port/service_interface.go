package port // inbound

import (
	"context"
	entity "product-service-api/internal/product/domain"
)

type ProductCommandServiceInterface interface {
	CreateProduct(ctx context.Context, product entity.Product, imageBytes []byte, imageFilename string) (entity.Product, error)
	UpdateProductByID(ctx context.Context,id string, product entity.Product, imageBytes []byte, imageFilename string) (entity.Product, error)
	DeleteProductByID(ctx context.Context, id string) error
}

type ProductQueryServiceInterface interface {
	GetProductByID(ctx context.Context, id string) (entity.Product, error)
	GetAllProducts(ctx context.Context) ([]entity.Product, error)
}
