package grpc

import (
	"context"
	"fmt"
	"product-service-api/internal/product/application/port"
	mapping "product-service-api/internal/product/adapter/handler/grpc/pb"
	"product-service-api/internal/product/adapter/handler/grpc/pb"


	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type productQueryHandler struct {
	pb.UnimplementedProductQueryServiceServer
	productQueryService   port.ProductQueryServiceInterface
}

func NewProductQueryHandler(pqs port.ProductQueryServiceInterface) *productQueryHandler {
	return &productQueryHandler{
		productQueryService:   pqs,
	}
}

func (ph *productQueryHandler) GetProductByID(ctx context.Context, productRequest *pb.GetProductByIDRequest) (*pb.ProductResponse, error) {
	if productRequest.GetId() == "" {
		return nil, fmt.Errorf("product id is required")
	}

	product, err := ph.productQueryService.GetProductByID(productRequest.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	productResponse :=  mapping.ProductEntityToResponse(product)

	return productResponse, nil
}

func (ph *productQueryHandler) GetAllProducts(ctx context.Context, _ *pb.EmptyRequest) (*pb.ListProductResponse, error) {
	products, err := ph.productQueryService.GetAllProducts()
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	productResponses := mapping.ListProductEntityToResponse(products)

	return &pb.ListProductResponse{
		Products: productResponses,
	}, nil
}
