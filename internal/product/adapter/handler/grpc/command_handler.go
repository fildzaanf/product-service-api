package grpc // inbound

import (
	"context"
	"product-service-api/internal/product/application/port"
	mapping "product-service-api/internal/product/adapter/handler/grpc/pb"
	"product-service-api/internal/product/adapter/handler/grpc/pb"
	"product-service-api/pkg/constant"
	"product-service-api/pkg/middleware"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type productCommandHandler struct {
	pb.UnimplementedProductCommandServiceServer
	productCommandService port.ProductCommandServiceInterface
	productQueryService   port.ProductQueryServiceInterface
}

func NewProductCommandHandler(pcs port.ProductCommandServiceInterface, pqs port.ProductQueryServiceInterface) *productCommandHandler {
	return &productCommandHandler{
		productCommandService: pcs,
		productQueryService:   pqs,
	}
}

func (ph *productCommandHandler) CreateProduct(ctx context.Context, productRequest *pb.CreateProductRequest) (*pb.ProductResponse, error) {

	userID, role, errExtract := middleware.ExtractTokenFromContext(ctx)
	if errExtract != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized access")
	}

	if role != constant.SELLER {
		return nil, status.Errorf(codes.PermissionDenied, constant.ERROR_ROLE_ACCESS)
	}

	imageBytes := productRequest.GetImageBytes()
	imageFilename := productRequest.GetImageFilename()

	productEntity := mapping.CreateProductRequestToEntity(productRequest, userID)

	createdProduct, err := ph.productCommandService.CreateProduct(productEntity, imageBytes, imageFilename)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	productResponse := mapping.ProductEntityToResponse(createdProduct)

	return productResponse, nil
}

func (ph *productCommandHandler) UpdateProduct(ctx context.Context, productRequest *pb.UpdateProductRequest) (*pb.ProductResponse, error) {

	userID, role, errExtract := middleware.ExtractTokenFromContext(ctx)
	if errExtract != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized access")
	}

	if productRequest.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "product id is required")
	}

	product, errGetProduct := ph.productQueryService.GetProductByID(productRequest.GetId())
	if errGetProduct != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	if product.UserID != userID {
		return nil, status.Error(codes.PermissionDenied, "forbidden access")
	}
	if role != constant.SELLER {
		return nil, status.Error(codes.PermissionDenied, constant.ERROR_ROLE_ACCESS)
	}

	imageBytes := productRequest.GetImageBytes()
	imageFilename := productRequest.GetImageFilename()
	productID := productRequest.GetId()

	productEntity := mapping.UpdateProductRequestToEntity(productRequest)

	updatedProduct, err := ph.productCommandService.UpdateProductByID(productID, productEntity, imageBytes, imageFilename)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	productResponse := mapping.ProductEntityToResponse(updatedProduct)

	return productResponse, nil
}

func (ph *productCommandHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {

	userID, role, errExtract := middleware.ExtractTokenFromContext(ctx)
	if errExtract != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized access")
	}

	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "product id is required")
	}

	product, err := ph.productQueryService.GetProductByID(req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	if product.UserID != userID {
		return nil, status.Error(codes.PermissionDenied, "forbidden access")
	}
	if role != constant.SELLER {
		return nil, status.Error(codes.PermissionDenied, constant.ERROR_ROLE_ACCESS)
	}

	if errDelete := ph.productCommandService.DeleteProductByID(req.GetId()); errDelete != nil {
		return nil, status.Error(codes.InvalidArgument, errDelete.Error())
	}

	return &mapping.DeleteProductResponse{}, nil
}
