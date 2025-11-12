package service

import (
	"context"
	"errors"
	"fmt"
	"product-service-api/infrastructure/cloud"
	"product-service-api/internal/product/application/port"
	entity "product-service-api/internal/product/domain"
	"product-service-api/pkg/constant"
	"product-service-api/pkg/validator"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type productCommandService struct {
	productCommandRepository port.ProductCommandRepositoryInterface
	productQueryRepository   port.ProductQueryRepositoryInterface
	userQueryClient          port.UserQueryClientInterface
}

func NewProductCommandService(pcr port.ProductCommandRepositoryInterface, pqr port.ProductQueryRepositoryInterface, uqr port.UserQueryClientInterface) port.ProductCommandServiceInterface {
	return &productCommandService{
		productCommandRepository: pcr,
		productQueryRepository:   pqr,
		userQueryClient:          uqr,
	}
}

func (pcs *productCommandService) CreateProduct(ctx context.Context, product entity.Product, imageBytes []byte, imageFilename string) (entity.Product, error) {

	_, errGetUser := pcs.userQueryClient.GetUserByID(ctx, product.UserID)
	if errGetUser != nil {
		return entity.Product{}, fmt.Errorf("invalid user: %w", errGetUser)
	}

	if len(imageBytes) > 0 {
		imageURL, errUpload := cloud.UploadImageBytesToS3(imageBytes, imageFilename)
		if errUpload != nil {
			return entity.Product{}, errUpload
		}
		product.ImageURL = imageURL
	}

	if err := validator.IsDataEmpty(
		[]string{"name", "description", "price", "stock"},
		product.Name, product.Description, product.Price, product.Stock,
	); err != nil {
		return entity.Product{}, err
	}

	if product.Price.LessThanOrEqual(decimal.NewFromInt(0)) {
		return entity.Product{}, errors.New(constant.ERROR_INVALID_PRICE)
	}

	if product.Stock < 0 {
		return entity.Product{}, errors.New(constant.ERROR_INVALID_STOCK)
	}

	if product.ID == "" {
		product.ID = uuid.New().String()
	}

	createdProduct, err := pcs.productCommandRepository.CreateProduct(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}

	return createdProduct, nil
}

func (pcs *productCommandService) UpdateProductByID(ctx context.Context, id string, product entity.Product, imageBytes []byte, imageFilename string) (entity.Product, error) {
	_, errGetUser := pcs.userQueryClient.GetUserByID(ctx, product.UserID)
	if errGetUser != nil {
		return entity.Product{}, fmt.Errorf("invalid user: %w", errGetUser)
	}

	existingProduct, err := pcs.productQueryRepository.GetProductByID(ctx, id)
	if err != nil {
		return entity.Product{}, errors.New(constant.ERROR_PRODUCT_NOT_FOUND)
	}

	if product.Name != "" && product.Name != existingProduct.Name {
		existingProduct.Name = product.Name
	}

	if product.Description != "" && product.Description != existingProduct.Description {
		existingProduct.Description = product.Description
	}
	if product.Price.GreaterThan(decimal.NewFromInt(0)) {
		existingProduct.Price = product.Price
	}
	if product.Stock >= 0 {
		existingProduct.Stock = product.Stock
	}

	if len(imageBytes) > 0 {
		imageURL, errUpload := cloud.UploadImageBytesToS3(imageBytes, imageFilename)
		if errUpload != nil {
			return entity.Product{}, errUpload
		}
		existingProduct.ImageURL = imageURL
	}

	existingProduct.ID = id

	updatedProduct, err := pcs.productCommandRepository.UpdateProductByID(ctx, id, existingProduct)
	if err != nil {
		return entity.Product{}, err
	}

	return updatedProduct, nil
}

func (pcs *productCommandService) DeleteProductByID(ctx context.Context, id string) error {

	_, err := pcs.productQueryRepository.GetProductByID(ctx, id)
	if err != nil {
		return errors.New(constant.ERROR_PRODUCT_NOT_FOUND)
	}

	if err := pcs.productCommandRepository.DeleteProductByID(ctx, id); err != nil {
		return err
	}

	return nil
}
