package rest // inbound

import (
	"net/http"
	"product-service-api/internal/product/application/port"
	"product-service-api/pkg/constant"
	"product-service-api/pkg/response"

	"github.com/labstack/echo/v4"
)

type productQueryHandler struct {
	productCommandService port.ProductCommandServiceInterface
	productQueryService   port.ProductQueryServiceInterface
}

func NewProductQueryHandler(pcs port.ProductCommandServiceInterface, pqs port.ProductQueryServiceInterface) *productQueryHandler {
	return &productQueryHandler{
		productCommandService: pcs,
		productQueryService:   pqs,
	}
}

// query
func (ph *productQueryHandler) GetProductByID(c echo.Context) error {
	ctx := c.Request().Context()
	
	productID := c.Param("id")
	if productID == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("product id is required"))
	}

	product, err := ph.productQueryService.GetProductByID(ctx, productID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse("product not found"))
	}

	productResponse := ProductEntityToResponse(product)
	return c.JSON(http.StatusOK, response.SuccessResponse(constant.SUCCESS_RETRIEVED, productResponse))
}

func (ph *productQueryHandler) GetAllProducts(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := ph.productQueryService.GetAllProducts(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to retrieve products"))
	}

	productResponses := ListProductEntityToResponse(products)

	return c.JSON(http.StatusOK, response.SuccessResponse(constant.SUCCESS_RETRIEVED, productResponses))
}
