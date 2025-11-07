package port // outbond

import (
	"context"
	entity "product-service-api/internal/product/domain"
)

type ProductCommandRepositoryInterface interface {
	CreateProduct(ctx context.Context, product entity.Product) (entity.Product, error)
	UpdateProductByID(ctx context.Context, id string, product entity.Product) (entity.Product, error)
	DeleteProductByID(ctx context.Context, id string) error
	UpdateProductStockByID(ctx context.Context, productID string, newStock int) error
}

type ProductQueryRepositoryInterface interface {
	GetProductByID(ctx context.Context, id string) (entity.Product, error)
	GetAllProducts(ctx context.Context) ([]entity.Product, error)
}
