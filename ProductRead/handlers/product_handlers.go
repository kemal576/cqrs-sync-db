package handlers

import (
	"ProductRead/repositories"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(repo *repositories.ProductReadRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := repo.GetAll(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

func GetProductById(repo *repositories.ProductReadRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		product, err := repo.GetById(context.Background(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}
