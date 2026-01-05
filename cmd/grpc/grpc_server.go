package main

import (
	"fmt"
	"log"
	"net"
	"product-service-api/infrastructure/config"
	"product-service-api/infrastructure/database"
	"product-service-api/pkg/middleware"

	"github.com/common-nighthawk/go-figure"

	productGRPC "product-service-api/internal/product/adapter/handler/grpc"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	middleware.InitLogger()

	godotenv.Load()
	config, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("[ERROR] failed to load configuration: %v", err)
	}

	psql := database.ConnectPostgreSQL(false)

	userServiceAddress := config.USERSERVICE.USER_GRPC_HOST + ":" + config.USERSERVICE.USER_GRPC_PORT
	userConn, err := grpc.NewClient(userServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("[ERROR] failed to connect to user service: %v", err)
	}
	defer userConn.Close()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.JWTUnaryInterceptor()),
	)

	productGRPC.RegisterProductServices(server, psql, userConn)

	reflection.Register(server)

	address := config.PRODUCTSERVICE.PRODUCT_GRPC_HOST + ":" + config.PRODUCTSERVICE.PRODUCT_GRPC_PORT
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("[FATAL] failed to listen product service: %v", err)
	}

	fig := figure.NewFigure("PRODUCT SERVICE API", "small", true)
	fig.Print()

	fmt.Printf("\n📡 Listening on %s\n", address)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
