package database

import (
	"fmt"
	"log"
	"product-service-api/infrastructure/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQL() *gorm.DB {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load MySQL configuration: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MYSQL.MYSQL_USER,
		config.MYSQL.MYSQL_PASS,
		config.MYSQL.MYSQL_HOST,
		config.MYSQL.MYSQL_PORT,
		config.MYSQL.MYSQL_NAME,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect MySQL: %v", err)
	}

	Migration(db)

	logrus.Info("connected to MySQL")

	return db
}
