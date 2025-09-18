package model

import (
	// "product-service-api/internal/user/adapter/model"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID          string          `gorm:"type:varchar(36);primaryKey"`
	UserID      string          `gorm:"type:varchar(36);not null"`
	Name        string          `gorm:"type:varchar(255);not null"`
	Description string          `gorm:"type:text"`
	Price       decimal.Decimal `gorm:"type:decimal(15,2);default:0.00;not null"`
	Stock       int             `gorm:"not null"`
	ImageURL    string          `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New().String()
	return nil
}
