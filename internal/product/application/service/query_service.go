package service

import (
	"context"
	"errors"
	"product-service-api/internal/product/application/port"
	entity "product-service-api/internal/product/domain"
	"product-service-api/pkg/constant"
)

type productQueryService struct {
	productQueryRepository   port.ProductQueryRepositoryInterface
	productCommandRepository port.ProductCommandRepositoryInterface
}

func NewProductQueryService(pqr port.ProductQueryRepositoryInterface, pcr port.ProductCommandRepositoryInterface) port.ProductQueryServiceInterface {
	return &productQueryService{
		productQueryRepository:   pqr,
		productCommandRepository: pcr,
	}
}

func (pqs *productQueryService) GetProductByID(ctx context.Context, id string) (entity.Product, error) {
	if id == "" {
		return entity.Product{}, errors.New(constant.ERROR_ID_INVALID)
	}

	product, err := pqs.productQueryRepository.GetProductByID(ctx, id)
	if err != nil {
		return entity.Product{}, errors.New("product not found")
	}

	return product, nil
}

func (pqs *productQueryService) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	products, err := pqs.productQueryRepository.GetAllProducts(ctx)

	if err != nil {
		return nil, errors.New(constant.ERROR_DATA_EMPTY)
	}

	return products, nil
}
