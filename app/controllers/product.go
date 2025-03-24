package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Product godoc
// @Summary Product
// @Description Product
// @Tags Product
// @Accept json
// @Produce json
// @Success 200
// @Router /v1/product [get]
func Product(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product Page",
	})
}
