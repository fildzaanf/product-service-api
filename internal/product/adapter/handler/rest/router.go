package rest

import (
	gormRepository "product-service-api/internal/product/adapter/repository/gorm"
	"product-service-api/internal/product/application/port"
	"product-service-api/internal/product/application/service"

	"product-service-api/pkg/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ProductRouter(product *echo.Group, db *gorm.DB, userQueryClient port.UserQueryClientInterface) {
	productQueryRepository := gormRepository.NewProductQueryRepository(db)
	productCommandRepository := gormRepository.NewProductCommandRepository(db)

	productQueryService := service.NewProductQueryService(productQueryRepository, productCommandRepository)
	productCommandService := service.NewProductCommandService(productCommandRepository, productQueryRepository, userQueryClient)

	productQueryHandler := NewProductQueryHandler(productCommandService, productQueryService)
	productCommandHandler := NewProductCommandHandler(productCommandService, productQueryService)

	product.POST("", productCommandHandler.CreateProduct, middleware.JWTMiddleware())
	product.PUT("/:id", productCommandHandler.UpdateProductByID, middleware.JWTMiddleware())
	product.DELETE("/:id", productCommandHandler.DeleteProductByID, middleware.JWTMiddleware())
	product.GET("/:id", productQueryHandler.GetProductByID)
	product.GET("", productQueryHandler.GetAllProducts)
}
