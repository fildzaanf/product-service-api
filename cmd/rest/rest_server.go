package main

import (
	"context"
	"product-service-api/infrastructure/config"
	"product-service-api/infrastructure/database"
	router "product-service-api/internal/product/adapter/handler/rest"
	"product-service-api/pkg/middleware"
	userClient "product-service-api/internal/product/adapter/client"

	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB, userHost, userPort string) {
	product := e.Group("/products")

	userQueryClient := userClient.NewUserRESTClient(userHost, userPort)

	router.ProductRouter(product, db, userQueryClient)
}

func main() {
	middleware.InitLogger()

	godotenv.Load()
	config, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("[ERROR] failed to load configuration: %v", err)
	}

	psql := database.ConnectPostgreSQL(true)

	e := echo.New()

	middleware.RemoveTrailingSlash(e)
	e.Use(middleware.Logger)
	middleware.RateLimiter(e)
	middleware.Recover(e)
	middleware.CORS(e)

	SetupRoutes(e, psql, config.USERSERVICE.USER_REST_HOST, config.USERSERVICE.USER_REST_PORT)

	address := config.PRODUCTSERVICE.PRODUCT_REST_HOST + ":" + config.PRODUCTSERVICE.PRODUCT_REST_PORT

	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		logrus.Info("[INFO] server is running on address ", address)
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit
		logrus.Warn("[WARNING] shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			logrus.Errorf("[ERROR] error shutting down server: %v", err)
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		logrus.Fatalf("[CRITICAL] server error: %v", err)
	case <-time.After(1 * time.Second):
		logrus.Info("[INFO] server is running smoothly...")
	}

	wg.Wait()
	logrus.Info("[INFO] server has been shut down gracefully.")
}
