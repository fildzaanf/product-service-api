package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	MYSQL          MySQLConfig
	POSTGRESQL     PostgreSQLConfig
	JWT            JWTConfig
	CLOUDSTORAGE   CloudStorageConfig
	MIDTRANS       MidtransConfig
	SMTP           SMTPConfig
	PRODUCTSERVICE ProductServiceConfig
	USERSERVICE    UserServiceConfig
}

type (
	MySQLConfig struct {
		MYSQL_USER string
		MYSQL_PASS string
		MYSQL_HOST string
		MYSQL_PORT string
		MYSQL_NAME string
	}

	PostgreSQLConfig struct {
		POSTGRESQL_USER string
		POSTGRESQL_PASS string
		POSTGRESQL_HOST string
		POSTGRESQL_PORT string
		POSTGRESQL_NAME string
	}

	JWTConfig struct {
		JWT_SECRET string
	}

	CloudStorageConfig struct {
		AWS_ACCESS_KEY_ID     string
		AWS_SECRET_ACCESS_KEY string
		AWS_REGION            string
		AWS_BUCKET_NAME       string
	}

	MidtransConfig struct {
		MIDTRANS_SERVER_KEY string
		MIDTRANS_CLIENT_KEY string
	}

	SMTPConfig struct {
		SMTP_USER string
		SMTP_PASS string
		SMTP_PORT string
		SMTP_HOST string
	}
	ProductServiceConfig struct {
		PRODUCT_REST_HOST string
		PRODUCT_REST_PORT string
		PRODUCT_GRPC_HOST string
		PRODUCT_GRPC_PORT string
	}
	UserServiceConfig struct {
		USER_REST_HOST string
		USER_REST_PORT string
		USER_GRPC_HOST string
		USER_GRPC_PORT string
	}
)

func LoadConfig() (*Configuration, error) {

	_, err := os.Stat(".env")
	if err == nil {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("failed to load environment variables from .env file: %w", err)
		}
	} else {
		fmt.Println("warning: .env file not found. make sure environment variables are set")
	}

	return &Configuration{
		MYSQL: MySQLConfig{
			MYSQL_USER: os.Getenv("MYSQL_USER"),
			MYSQL_PASS: os.Getenv("MYSQL_PASS"),
			MYSQL_HOST: os.Getenv("MYSQL_HOST"),
			MYSQL_PORT: os.Getenv("MYSQL_PORT"),
			MYSQL_NAME: os.Getenv("MYSQL_NAME"),
		},
		POSTGRESQL: PostgreSQLConfig{
			POSTGRESQL_USER: os.Getenv("POSTGRESQL_USER"),
			POSTGRESQL_PASS: os.Getenv("POSTGRESQL_PASS"),
			POSTGRESQL_HOST: os.Getenv("POSTGRESQL_HOST"),
			POSTGRESQL_PORT: os.Getenv("POSTGRESQL_PORT"),
			POSTGRESQL_NAME: os.Getenv("POSTGRESQL_NAME"),
		},
		JWT: JWTConfig{
			JWT_SECRET: os.Getenv("JWT_SECRET"),
		},
		CLOUDSTORAGE: CloudStorageConfig{
			AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			AWS_REGION:            os.Getenv("AWS_REGION"),
			AWS_BUCKET_NAME:       os.Getenv("AWS_BUCKET_NAME"),
		},
		MIDTRANS: MidtransConfig{
			MIDTRANS_SERVER_KEY: os.Getenv("MIDTRANS_SERVER_KEY"),
			MIDTRANS_CLIENT_KEY: os.Getenv("MIDTRANS_CLIENT_KEY"),
		},
		SMTP: SMTPConfig{
			SMTP_USER: os.Getenv("SMTP_USER"),
			SMTP_PASS: os.Getenv("SMTP_PASS"),
			SMTP_PORT: os.Getenv("SMTP_PORT"),
			SMTP_HOST: os.Getenv("SMTP_HOST"),
		},
		PRODUCTSERVICE: ProductServiceConfig{
			PRODUCT_REST_HOST: os.Getenv("PRODUCT_REST_HOST"),
			PRODUCT_REST_PORT: os.Getenv("PRODUCT_REST_PORT"),
			PRODUCT_GRPC_HOST: os.Getenv("PRODUCT_GRPC_HOST"),
			PRODUCT_GRPC_PORT: os.Getenv("PRODUCT_GRPC_PORT"),
		},
		USERSERVICE: UserServiceConfig{	
			USER_REST_HOST: os.Getenv("USER_REST_HOST"),
			USER_REST_PORT: os.Getenv("USER_REST_PORT"),
			USER_GRPC_HOST: os.Getenv("USER_GRPC_HOST"),
			USER_GRPC_PORT: os.Getenv("USER_GRPC_PORT"),
		},
	}, nil
}
