package grpc

import (
	"product-service-api/internal/product/application/service"
	gormRepository "product-service-api/internal/product/adapter/repository/gorm"
	"product-service-api/internal/product/adapter/handler/grpc/pb"
	userClient "product-service-api/internal/product/adapter/client"

	grpc "google.golang.org/grpc"
	"gorm.io/gorm"
)

func RegisterProductServices(server *grpc.Server, db *gorm.DB, conn *grpc.ClientConn) {
	productQueryRepository := gormRepository.NewProductQueryRepository(db)
	productCommandRepository := gormRepository.NewProductCommandRepository(db)
	
	userQueryClient := userClient.NewUserGRPCClient(conn)

	productQueryService := service.NewProductQueryService(productQueryRepository, productCommandRepository)
	productCommandService := service.NewProductCommandService(productCommandRepository, productQueryRepository, userQueryClient)

	productQueryHandler := NewProductQueryHandler(productQueryService)
	productCommandHandler := NewProductCommandHandler(productCommandService, productQueryService)

	pb.RegisterProductQueryServiceServer(server, productQueryHandler)
	pb.RegisterProductCommandServiceServer(server, productCommandHandler)
}
