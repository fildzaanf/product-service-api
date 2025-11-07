package pb // inbound

import (
	entity "product-service-api/internal/product/domain"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateProductRequestToEntity(request *CreateProductRequest, userID string) entity.Product {
	price, _ := decimal.NewFromString(request.Price)
	return entity.Product{
		UserID:      userID,
		Name:        request.Name,
		Description: request.Description,
		Price:       price,
		Stock:       int(request.Stock),
	}
}

func UpdateProductRequestToEntity(request *UpdateProductRequest) entity.Product {

	product := entity.Product{
		UserID: request.UserId,
	}

	if request.Name != "" {
		product.Name = request.Name
	}
	if request.Description != "" {
		product.Description = request.Description
	}

	if request.Price != "" {
		price, err := decimal.NewFromString(request.Price)
		if err != nil {
			return entity.Product{}
		}
		product.Price = price
	}

	if request.Stock >= 0 {
		product.Stock = int(request.Stock)
	}

	return product
}

func ProductEntityToResponse(product entity.Product) *ProductResponse {
	return &ProductResponse{
		Id:          product.ID,
		UserId:      product.UserID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price.String(),
		Stock:       int32(product.Stock),
		ImageUrl:    product.ImageURL,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}

func ListProductEntityToResponse(products []entity.Product) []*ProductResponse {
	productResponses := make([]*ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = ProductEntityToResponse(product)
	}
	return productResponses
}
