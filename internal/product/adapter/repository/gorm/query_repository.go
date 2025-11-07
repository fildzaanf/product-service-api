package gorm

import (
	"context"
	"errors"
	"product-service-api/internal/product/adapter/model"
	"product-service-api/internal/product/application/port"
	entity "product-service-api/internal/product/domain"

	"gorm.io/gorm"
)

type productQueryRepository struct {
	db *gorm.DB
}

func NewProductQueryRepository(db *gorm.DB) port.ProductQueryRepositoryInterface {
	return &productQueryRepository{
		db: db,
	}
}

func (pr *productQueryRepository) GetProductByID(ctx context.Context, id string) (entity.Product, error) {
	var product model.Product
	result := pr.db.Where("id = ?", id).First(&product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return entity.Product{}, errors.New("product not found")
		}
		return entity.Product{}, result.Error
	}

	return entity.ProductModelToEntity(product), nil
}

func (pr *productQueryRepository) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	var products []model.Product
	result := pr.db.Find(&products)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("no products found")
		}
		return nil, result.Error
	}

	return entity.ListProductModelToEntity(products), nil
}
