package database

import (
	"fmt"
	"product-service-api/infrastructure/config"

	"log"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectPostgreSQL(migrate bool) *gorm.DB {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load PostgreSQL configuration: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC prefer_simple_protocol=false disable_automatic_ping=true",
		cfg.POSTGRESQL.POSTGRESQL_HOST,
		cfg.POSTGRESQL.POSTGRESQL_USER,
		cfg.POSTGRESQL.POSTGRESQL_PASS,
		cfg.POSTGRESQL.POSTGRESQL_NAME,
		cfg.POSTGRESQL.POSTGRESQL_PORT,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect PostgreSQL: %v", err)
	}

	if migrate {
		Migration(db)
	}

	logrus.Info("[INFO] connected to PostgreSQL")

	return db
}
