package handlers

import (
	"oracle-demo/services"

	"github.com/gin-gonic/gin"
)

func GenerateC0401Handler(c *gin.Context) {
	// 處理創建發票的邏輯
	segment_no := c.Param("segment_no")
	

	// Call the handler to generate the invoice
	invoice_data, err := services.GenerateC0401Invoice(segment_no)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate invoice"})
		return
	}
	c.JSON(200, gin.H{"data": invoice_data})
}
