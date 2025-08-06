package handlers

import (
	"oracle-demo/services"
	"time"

	"github.com/gin-gonic/gin"
)

func GenerateC0401Handler(c *gin.Context) {
	// 處理創建發票的邏輯
	segment_no := c.Param("segment_no")
	invoice_date := c.Query("invoice_date")
	if invoice_date == "" {
		invoice_date = time.Now().Format("20060102")
	}

	// Call the handler to generate the invoice
	invoice_data, err := services.GenerateC0401Invoice(segment_no, invoice_date)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate invoice"})
		return
	}
	c.JSON(200, gin.H{"data": invoice_data})
}
