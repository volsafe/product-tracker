package handlers

import (
	"net/http"
	"product-tracker/storage"

	"github.com/gin-gonic/gin"
)

// Product represents the product request/response structure
// @Description Product information
type Product struct {
	Name           string  `json:"name" example:"Product A" binding:"required"`
	Quantity       int     `json:"quantity" example:"100" binding:"required,min=0"`
	EnergyConsumed float64 `json:"energy_consumed" example:"50.5" binding:"required,min=0"`
	Date           string  `json:"date" example:"2024-03-10" binding:"required,datetime=2006-01-02"`
}

// ImportProduct godoc
// @Summary      Import a new product
// @Description  Import a single product with its details
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      Product  true  "Product object"
// @Success      201      {object}  Product
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /product/insert [post]
// @Security     BearerAuth
func ImportProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storageInstance, err := storage.NewStorage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer storageInstance.Close()

	if err := storageInstance.InsertProduct(c.Request.Context(), storage.Product(product)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProducts godoc
// @Summary      List all products
// @Description  Get a list of all products in the system
// @Tags         products
// @Produce      json
// @Success      200  {array}   Product
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /product/list [get]
// @Security     BearerAuth
func GetProducts(c *gin.Context) {
	storageInstance, err := storage.NewStorage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer storageInstance.Close()

	products, err := storageInstance.GetProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductsByName godoc
// @Summary      Get products by name
// @Description  Get a list of products filtered by name
// @Tags         products
// @Produce      json
// @Param        name  path      string  true  "Product name"
// @Success      200   {array}   Product
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /product/list/{name} [get]
// @Security     BearerAuth
func GetProductsByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
		return
	}

	storageInstance, err := storage.NewStorage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer storageInstance.Close()

	products, err := storageInstance.GetProductsByName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
