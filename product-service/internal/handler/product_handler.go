package handler

import (
	"net/http"

	"product-service/internal/form"
	"product-service/internal/service"
	"product-service/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var form form.CreateProductForm

	if err := c.ShouldBindJSON(&form); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := util.FormatValidationError(validationErrors)
			util.RespondWithError(c, http.StatusBadRequest, "Validation failed", formattedErrors)
			return
		}
		util.RespondWithError(c, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	err := h.ProductService.CreateProduct(form.Name, form.Description, form.Price, form.Stock) // We will update category later on
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Unable to create product", nil)
		return
	}
	util.RespondWithSuccess(c, "product created", nil)
	return
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	searchQuery, _ := c.GetQuery("search")
	products, err := h.ProductService.GetProducts(searchQuery)
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Unable to fetch the products", nil)
		return
	}
	util.RespondWithSuccess(c, "products fetched", products)
	return
}

func (h *ProductHandler) GetProductById(c *gin.Context) {
	id := c.Param("id")
	product, err := h.ProductService.GetProductById(id)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Unable to fetch the product", nil)
		return
	}
	util.RespondWithSuccess(c, "products fetched", product)
	return
}

func (h *ProductHandler) UpdateProductById(c *gin.Context) {
	id := c.Param("id")

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	updatedProduct, err := h.ProductService.UpdateProductById(id, updateData)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Unable to fetch the product", nil)
		return
	}
	util.RespondWithSuccess(c, "products fetched", updatedProduct)
	return
}

func (h *ProductHandler) DeleteProductById(c *gin.Context) {
	id := c.Param("id")

	err := h.ProductService.DeleteProductById(id)
	if err != nil {
		util.RespondWithError(c, http.StatusNotFound, "Unable to delete product", nil)
		return
	}

	util.RespondWithSuccess(c, "product deleted", nil)
}
