package database

import (
	"log"

	"product-service-api/internal/product/adapter/model"

	"gorm.io/gorm"
)

func Migration(db *gorm.DB) {
	migrator := db.Migrator()

	db.AutoMigrate(
		&model.Product{},
	)

	tables := []string{"products"}
	for _, table := range tables {
		if !migrator.HasTable(table) {
			log.Fatalf("table %s was not successfully created", table)
		}
	}
	log.Println("all tables were successfully migrated")
}
