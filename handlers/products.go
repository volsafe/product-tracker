package handlers

import (
	"fmt"
	"io"
	"product-tracker/controllers"

	"github.com/gin-gonic/gin"
	"encoding/json"
)


type Product struct {
	Name  string `json:"name"`
	Quantity int `json:"quantity"`
	EnergyConsumed float64 `json:"energy_consumed"`
	Date string `json:"date"`
}

func ImportProduct(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error(),})
		return
	}
	defer c.Request.Body.Close()

	// Unmarshal the JSON payload
	var p Product
	err = json.Unmarshal(body, &p)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error(),})
		fmt.Printf("Error: %v", err)
		return
	}
	// Call the controller
	np := controllers.Product{
		Name: p.Name,
		Quantity: p.Quantity,
		EnergyConsumed: p.EnergyConsumed,
		Date: p.Date,
	}
	err = controllers.InsertProduct(c, np)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error(),})
		return
	}

}

func GetProductsByName(c *gin.Context) {
	name := c.Param("name")
	products, err := controllers.GetProductsByName(c, name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error(),})
		return
	}
	c.JSON(200, products)
}

func GetProducts(c *gin.Context) {
	products, err := controllers.GetProducts(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error(),})
		return
	}
	c.JSON(200, products)
}


