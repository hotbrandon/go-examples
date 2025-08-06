package services

import (
	"log"
	"oracle-demo/models"
)

func GenerateC0401Invoice(segment_no string) ([]models.Invoice, error) {
	// 從 Oracle 產出已開立發票
	log.Println("Generating C0401 invoice for segment:", segment_no)
	db, err := models.GetOracleConnection(segment_no)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	invoice_data := []models.Invoice{
		{InvoiceID: "1"},
		{InvoiceID: "2"},
		{InvoiceID: "3"},
	}

	return invoice_data, nil
}
