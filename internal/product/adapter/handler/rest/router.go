package rest


import (
	"product-service-api/internal/product/application/service"
	gormRepository "product-service-api/internal/product/adapter/repository/gorm"

	"product-service-api/pkg/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ProductRouter(product *echo.Group, db *gorm.DB) {
	productQueryRepository := gormRepository.NewProductQueryRepository(db)
	productCommandRepository := gormRepository.NewProductCommandRepository(db)

	productQueryService := service.NewProductQueryService(productQueryRepository, productCommandRepository)
	productCommandService := service.NewProductCommandService(productCommandRepository, productQueryRepository, nil)

	productQueryHandler := NewProductQueryHandler(productCommandService, productQueryService)
	productCommandHandler := NewProductCommandHandler(productCommandService, productQueryService)

	product.POST("", productCommandHandler.CreateProduct, middleware.JWTMiddleware())
	product.PUT("/:id", productCommandHandler.UpdateProductByID, middleware.JWTMiddleware())
	product.DELETE("/:id", productCommandHandler.DeleteProductByID, middleware.JWTMiddleware())
	product.GET("/:id", productQueryHandler.GetProductByID)
	product.GET("", productQueryHandler.GetAllProducts)
}