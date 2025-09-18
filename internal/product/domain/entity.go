package domain


import (
	"product-service-api/internal/product/adapter/model"
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID          string
	UserID      string
	Name        string
	Description string
	Price       decimal.Decimal
	Stock       int
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// mapper
func ProductEntityToModel(productEntity Product) model.Product {
	return model.Product{
		ID:          productEntity.ID,
		UserID:      productEntity.UserID,
		Name:        productEntity.Name,
		Description: productEntity.Description,
		Price:       productEntity.Price,
		Stock:       productEntity.Stock,
		ImageURL:    productEntity.ImageURL,
		CreatedAt:   productEntity.CreatedAt,
		UpdatedAt:   productEntity.UpdatedAt,
		DeletedAt:   productEntity.DeletedAt,
	}
}

func ProductModelToEntity(productModel model.Product) Product {
	return Product{
		ID:          productModel.ID,
		UserID:      productModel.UserID,
		Name:        productModel.Name,
		Description: productModel.Description,
		Price:       productModel.Price,
		Stock:       productModel.Stock,
		ImageURL:    productModel.ImageURL,
		CreatedAt:   productModel.CreatedAt,
		UpdatedAt:   productModel.UpdatedAt,
		DeletedAt:   productModel.DeletedAt,
	}
}

func ListProductEntityToModel(productEntities []Product) []model.Product {
	listProductModels := []model.Product{}
	for _, product := range productEntities {
		listProductModels = append(listProductModels, ProductEntityToModel(product))
	}
	return listProductModels
}

func ListProductModelToEntity(productModels []model.Product) []Product {
	listProductEntities := []Product{}
	for _, product := range productModels {
		listProductEntities = append(listProductEntities, ProductModelToEntity(product))
	}
	return listProductEntities
}
