package grpc

import (
	"log"
	"net"
	"product-service-api/infrastructure/config"
	"product-service-api/infrastructure/database"
	"product-service-api/pkg/middleware"

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

	userServiceAddress := config.USERSERVICE.USER_HOST + ":" + config.USERSERVICE.USER_PORT
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

	address := config.PRODUCTSERVICE.PRODUCT_HOST + ":" + config.PRODUCTSERVICE.PRODUCT_PORT
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("[FATAL] failed to listen product service: %v", err)
	}

	log.Printf("[INFO] product service running on %s", address)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("[FATAL] failed to serve product service: %v", err)
	}
}
