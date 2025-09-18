package rest // inbound

import (
	entity "product-service-api/internal/product/domain"
	"time"

	"github.com/shopspring/decimal"
)

// request
type (
	CreateProductRequest struct {
		Name        string          `json:"name" form:"name"`
		Description string          `json:"description" form:"description"`
		Price       decimal.Decimal `json:"price" form:"price"`
		Stock       int             `json:"stock" form:"stock"`
		ImageURL    string          `json:"image_url" form:"image_url"`
	}

	UpdateProductRequest struct {
		Name        string          `json:"name" form:"name"`
		Description string          `json:"description" form:"description"`
		Price       decimal.Decimal `json:"price" form:"price"`
		Stock       int             `json:"stock" form:"stock"`
		ImageURL    string          `json:"image_url" form:"image_url"`
	}
)

func CreateProductRequestToEntity(request CreateProductRequest, userID string) entity.Product {
	return entity.Product{
		UserID:      userID,
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
		ImageURL:    request.ImageURL,
	}
}

func UpdateProductRequestToEntity(request UpdateProductRequest) entity.Product {
	return entity.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
		ImageURL:    request.ImageURL,
	}
}

// response

type (
	ProductResponse struct {
		ID          string          `json:"id"`
		UserID      string          `json:"user_id"`
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Price       decimal.Decimal `json:"price"`
		Stock       int             `json:"stock"`
		ImageURL    string          `json:"image_url"`
		CreatedAt   time.Time       `json:"created_at"`
		UpdatedAt   time.Time       `json:"updated_at"`
	}
)

func ProductEntityToResponse(product entity.Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		UserID:      product.UserID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		ImageURL:    product.ImageURL,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func ListProductEntityToResponse(products []entity.Product) []ProductResponse {
	productResponses := make([]ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = ProductEntityToResponse(product)
	}
	return productResponses
}
