package gorm


import (
	"errors"
	"fmt"
	entity"product-service-api/internal/product/domain"
	"product-service-api/internal/product/application/port"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type productCommandRepository struct {
	db *gorm.DB
}

func NewProductCommandRepository(db *gorm.DB) port.ProductCommandRepositoryInterface {
	return &productCommandRepository{
		db: db,
	}
}

func (pcr *productCommandRepository) CreateProduct(product entity.Product) (entity.Product, error) {
	tx := pcr.db.Begin()
	if tx.Error != nil {
		return entity.Product{}, tx.Error
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return entity.Product{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (pcr *productCommandRepository) UpdateProductByID(id string, product entity.Product) (entity.Product, error) {
	tx := pcr.db.Begin()
	if tx.Error != nil {
		return entity.Product{}, tx.Error
	}

	existingProduct := entity.Product{}
	if err := tx.First(&existingProduct, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return entity.Product{}, errors.New("product not found")
	}

	updateFields := map[string]interface{}{}

	if product.Name != "" {
		updateFields["name"] = product.Name
	}
	if product.Description != "" {
		updateFields["description"] = product.Description
	}
	if product.Price.GreaterThan(decimal.NewFromInt(0)) {
		updateFields["price"] = product.Price
	}
	if product.Stock >= 0 {
		updateFields["stock"] = product.Stock
	}
	if product.ImageURL != "" {
		updateFields["image_url"] = product.ImageURL
	}


	if len(updateFields) > 0 {
		updateFields["updated_at"] = time.Now()

		if err := tx.Model(&existingProduct).Select("*").Updates(updateFields).Error; err != nil {
			tx.Rollback()
			return entity.Product{}, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Println("Commit Error:", err)
		return entity.Product{}, err
	}

	return existingProduct, nil
}

func (pcr *productCommandRepository) DeleteProductByID(id string) error {
	tx := pcr.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("id = ?", id).Delete(&entity.Product{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (pcr *productCommandRepository) UpdateProductStockByID(productID string, newStock int) error {
	tx := pcr.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Model(&entity.Product{}).
		Where("id = ?", productID).
		Update("stock", newStock).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
