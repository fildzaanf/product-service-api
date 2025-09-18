package port // outbond


import entity "product-service-api/internal/product/domain"

type ProductCommandRepositoryInterface interface {
	CreateProduct(product entity.Product) (entity.Product, error)
	UpdateProductByID(id string, product entity.Product) (entity.Product, error)
	DeleteProductByID(id string) error
	UpdateProductStockByID(productID string, newStock int) error
}

type ProductQueryRepositoryInterface interface {
	GetProductByID(id string) (entity.Product, error)
	GetAllProducts() ([]entity.Product, error)
}
