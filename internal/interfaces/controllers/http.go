package controllers

import (
	"fmt"
	"net/http"
	"sagio/internal/interfaces/dtos"

	"github.com/gin-gonic/gin"
)

func InitTransaction(c *gin.Context) {
	var data dtos.InitTransactionDTO
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"detail": fmt.Sprintf("Invalid request body: %v", err),
		})
	}
}
