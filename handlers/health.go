package handlers

import (
	"net/http"

	"product-tracker/controllers"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary      Health check endpoint
// @Description  Check if the API is up and running
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	err := controllers.HealthCheck(c)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "API is healthy",
	})
}
