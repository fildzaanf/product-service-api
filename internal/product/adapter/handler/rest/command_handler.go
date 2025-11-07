package rest // inbound

import (
	"io"
	"net/http"
	"product-service-api/internal/product/application/port"
	"product-service-api/pkg/constant"
	"product-service-api/pkg/middleware"
	"product-service-api/pkg/response"

	"github.com/labstack/echo/v4"
)

type productCommandHandler struct {
	productCommandService port.ProductCommandServiceInterface
	productQueryService   port.ProductQueryServiceInterface
}

func NewProductCommandHandler(pcs port.ProductCommandServiceInterface, pqs port.ProductQueryServiceInterface) *productCommandHandler {
	return &productCommandHandler{
		productCommandService: pcs,
		productQueryService:   pqs,
	}
}

func (h *productCommandHandler) CreateProduct(c echo.Context) error {
	 ctx := c.Request().Context()
	userID, role, errExtract := middleware.ExtractToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse("unauthorized access"))
	}

	if role != constant.SELLER {
		return c.JSON(http.StatusForbidden, response.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	var productRequest CreateProductRequest
	if err := c.Bind(&productRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
	}

	var imageBytes []byte
	var imageFilename string

	imageFile, err := c.FormFile("image_url")
	if imageFile != nil && err == nil {
		src, err := imageFile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse("failed to read image"))
		}
		defer src.Close()

		imageBytes, err = io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse("failed to read image bytes"))
		}
		imageFilename = imageFile.Filename
	}

	product := CreateProductRequestToEntity(productRequest, userID)

	createdProduct, err := h.productCommandService.CreateProduct(ctx, product, imageBytes, imageFilename)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
	}

	productResponse := ProductEntityToResponse(createdProduct)

	return c.JSON(http.StatusCreated, response.SuccessResponse(constant.SUCCESS_CREATED, productResponse))

}

func (ph *productCommandHandler) UpdateProductByID(c echo.Context) error {
	ctx := c.Request().Context()

	userID, role, errExtract := middleware.ExtractToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse("unauthorized access"))
	}

	productID := c.Param("id")
	if productID == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("product id is required"))
	}

	product, err := ph.productQueryService.GetProductByID(ctx, productID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse("product not found"))
	}

	if product.UserID != userID {
		return c.JSON(http.StatusForbidden, response.ErrorResponse("forbidden access"))
	}

	if role != constant.SELLER {
		return c.JSON(http.StatusForbidden, response.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	var productRequest UpdateProductRequest
	if err := c.Bind(&productRequest); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
	}

	var imageBytes []byte
	var imageFilename string

	imageFile, err := c.FormFile("image_url")
	if imageFile != nil && err == nil {
		src, err := imageFile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse("failed to read image"))
		}
		defer src.Close()

		imageBytes, err = io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse("failed to read image bytes"))
		}
		imageFilename = imageFile.Filename
	}

	productEntity := UpdateProductRequestToEntity(productRequest)

	updatedProduct, err := ph.productCommandService.UpdateProductByID(ctx, productID, productEntity, imageBytes, imageFilename)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
	}

	productResponse := ProductEntityToResponse(updatedProduct)
	return c.JSON(http.StatusOK, response.SuccessResponse(constant.SUCCESS_UPDATED, productResponse))
}

func (ph *productCommandHandler) DeleteProductByID(c echo.Context) error {
	ctx := c.Request().Context()

	userID, role, errExtract := middleware.ExtractToken(c)
	if errExtract != nil {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse("unauthorized access"))
	}

	productID := c.Param("id")
	if productID == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("product id is required"))
	}

	product, err := ph.productQueryService.GetProductByID(ctx, productID)
	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse("product not found"))
	}

	if product.UserID != userID {
		return c.JSON(http.StatusForbidden, response.ErrorResponse("forbidden access"))
	}

	if role != constant.SELLER {
		return c.JSON(http.StatusForbidden, response.ErrorResponse(constant.ERROR_ROLE_ACCESS))
	}

	if err := ph.productCommandService.DeleteProductByID(ctx, productID); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("product successfully deleted", nil))
}
