package service

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"oracle-demo/models"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/sijms/go-ora/v2"
)

// InvoiceRow represents a row from the database query
// Using sql.NullString to properly handle Oracle NULLs
type InvoiceRow struct {
	InvoiceNumber        sql.NullString
	InvoiceDate          sql.NullString
	InvoiceTime          sql.NullString
	BuyerIdentifier      sql.NullString
	BuyerName            sql.NullString
	BuyerAddress         sql.NullString
	BuyerTelephoneNumber sql.NullString
	BuyerEmailAddress    sql.NullString
	SalesAmount          sql.NullString
	FreeTaxSalesAmount   sql.NullString
	ZeroTaxSalesAmount   sql.NullString
	TaxType              sql.NullString
	TaxRate              sql.NullString
	TaxAmount            sql.NullString
	TotalAmount          sql.NullString
	PrintMark            sql.NullString
	RandomNumber         sql.NullString
	MainRemark           sql.NullString
	CarrierType          sql.NullString
	CarrierID1           sql.NullString
	CarrierID2           sql.NullString
	NPOBAN               sql.NullString
	LineNo               int
	Description          sql.NullString
	Quantity             sql.NullString
	UnitPrice            sql.NullString
	Amount               sql.NullString
	DetailTaxType        sql.NullString
	Remark               sql.NullString
}

// CSVGenerationResult represents the result of CSV generation
type CSVGenerationResult struct {
	InvoiceCount int
	TotalRows    int
	FilePath     string
}

// InvoiceService handles invoice-related operations
type InvoiceService struct {
	rootDir string
}

// NewInvoiceService creates a new InvoiceService
func NewInvoiceService(rootDir string) *InvoiceService {
	return &InvoiceService{
		rootDir: rootDir,
	}
}

// GenC0401 generates C0401 CSV file from invoice data
func (s *InvoiceService) GenC0401(segmentNo, invoiceDate string) (*CSVGenerationResult, error) {
	db, err := models.GetOracleConnection()
	if err != nil {
		return nil, fmt.Errorf("getting database connection: %w", err)
	}
	defer db.Close()

	// Parse invoice date
	parsedDate, err := time.Parse("20060102", invoiceDate)
	if err != nil {
		return nil, fmt.Errorf("parsing invoice date: %w", err)
	}

	// Call stored procedures
	if err := s.callStoredProcedures(db, segmentNo, parsedDate); err != nil {
		return nil, fmt.Errorf("calling stored procedures: %w", err)
	}

	// Execute main query
	rows, err := s.executeQuery(db)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	// Generate CSV file
	result, err := s.generateCSVFile(rows, segmentNo, invoiceDate)
	if err != nil {
		return nil, fmt.Errorf("generating CSV file: %w", err)
	}

	return result, nil
}

// callStoredProcedures calls the required Oracle stored procedures
func (s *InvoiceService) callStoredProcedures(db *sql.DB, segmentNo string, invoiceDate time.Time) error {
	// Call PK_ERP.P_SET_SEGMENT_NO
	_, err := db.Exec("BEGIN PK_ERP.P_SET_SEGMENT_NO(:1); END;", strings.ToUpper(segmentNo))
	if err != nil {
		return fmt.Errorf("calling PK_ERP.P_SET_SEGMENT_NO: %w", err)
	}

	// Call ARGOERP.P_RPT_MRIF004
	_, err = db.Exec("BEGIN ARGOERP.P_RPT_MRIF004(:1); END;", invoiceDate)
	if err != nil {
		return fmt.Errorf("calling ARGOERP.P_RPT_MRIF004: %w", err)
	}

	return nil
}

// executeQuery executes the main data retrieval query
// Using TRIM with proper NULL handling for Oracle
func (s *InvoiceService) executeQuery(db *sql.DB) ([]InvoiceRow, error) {
	query := `
		SELECT TRIM(var_attr01) InvoiceNumber,
			TRIM(DECODE(num_attr08,1,TO_CHAR(date_attr01,'YYYYMMDD'),'')) InvoiceDate,
			TRIM(DECODE(num_attr08,1,var_attr02,'')) InvoiceTime,        
			TRIM(DECODE(num_attr08,1,var_attr03,'')) BuyerIdentifier,    
			TRIM(DECODE(num_attr08,1,var_attr04,'')) BuyerName,          
			TRIM(DECODE(num_attr08,1,var_attr05,'')) BuyerAddress,       
			TRIM(DECODE(num_attr08,1,var_attr06,'')) BuyerTelephoneNumber,
			TRIM(DECODE(num_attr08,1,var_attr07,'')) BuyerEmailAddress,  
			TRIM(TO_CHAR(DECODE(num_attr08,1,num_attr01,NULL))) SalesAmount,        
			TRIM(TO_CHAR(DECODE(num_attr08,1,num_attr02,NULL))) FreeTaxSalesAmount, 
			TRIM(TO_CHAR(DECODE(num_attr08,1,num_attr03,NULL))) ZeroTaxSalesAmount, 
			TRIM(DECODE(num_attr08,1,var_attr08,'')) TaxType,            
			TRIM(DECODE(num_attr08,1,TO_CHAR(num_attr04,'0.99'),'')) TaxRate,
			TRIM(TO_CHAR(DECODE(num_attr08,1,num_attr05,NULL))) TaxAmount,   
			TRIM(TO_CHAR(DECODE(num_attr08,1,num_attr06,NULL))) TotalAmount, 
			TRIM(DECODE(num_attr08,1,var_attr09,'')) PrintMark,   
			TRIM(TO_CHAR(DECODE(num_attr08,1,num_attr07,NULL))) RandomNumber,
			TRIM(DECODE(num_attr08,1,var_attr10,'')) MainRemark,  
			TRIM(DECODE(num_attr08,1,var_attr11,'')) CarrierType, 
			TRIM(DECODE(num_attr08,1,var_attr12,'')) CarrierId1,  
			TRIM(DECODE(num_attr08,1,var_attr13,'')) CarrierId2,  
			TRIM(DECODE(num_attr08,1,var_attr14,'')) NPOBAN,      
			num_attr08 line_no,    
			TRIM(var_attr15) Description,
			TRIM(TO_CHAR(num_attr09)) Quantity,   
			TRIM(TO_CHAR(num_attr10)) UnitPrice,  
			TRIM(TO_CHAR(num_attr11)) Amount,
			TRIM(var_attr08) DetailTaxType,     
			TRIM(var_attr16) Remark      
		FROM argoerp.mr_global_temp a
		WHERE a.pid = 'MRIF004'
		ORDER BY a.var_attr01, a.num_attr08`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	var invoiceRows []InvoiceRow
	for rows.Next() {
		var row InvoiceRow
		err := rows.Scan(
			&row.InvoiceNumber,
			&row.InvoiceDate,
			&row.InvoiceTime,
			&row.BuyerIdentifier,
			&row.BuyerName,
			&row.BuyerAddress,
			&row.BuyerTelephoneNumber,
			&row.BuyerEmailAddress,
			&row.SalesAmount,
			&row.FreeTaxSalesAmount,
			&row.ZeroTaxSalesAmount,
			&row.TaxType,
			&row.TaxRate,
			&row.TaxAmount,
			&row.TotalAmount,
			&row.PrintMark,
			&row.RandomNumber,
			&row.MainRemark,
			&row.CarrierType,
			&row.CarrierID1,
			&row.CarrierID2,
			&row.NPOBAN,
			&row.LineNo,
			&row.Description,
			&row.Quantity,
			&row.UnitPrice,
			&row.Amount,
			&row.DetailTaxType,
			&row.Remark,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		invoiceRows = append(invoiceRows, row)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return invoiceRows, nil
}

// generateCSVFile creates the CSV file with invoice data
func (s *InvoiceService) generateCSVFile(rows []InvoiceRow, segmentNo, invoiceDate string) (*CSVGenerationResult, error) {
	// Create directory path
	csvDirPath := filepath.Join(s.rootDir, "cxnvol", segmentNo)
	if err := os.MkdirAll(csvDirPath, 0755); err != nil {
		return nil, fmt.Errorf("creating directory: %w", err)
	}

	// Create CSV file
	csvFileName := fmt.Sprintf("C0401-%s-%s.csv", invoiceDate, segmentNo)
	csvFilePath := filepath.Join(csvDirPath, csvFileName)

	file, err := os.Create(csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("creating CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"InvoiceNumber", "InvoiceDate", "InvoiceTime", "BuyerIdentifier", "BuyerName",
		"BuyerAddress", "BuyerTelephoneNumber", "BuyerEmailAddress", "SalesAmount",
		"FreeTaxSalesAmount", "ZeroTaxSalesAmount", "TaxType", "TaxRate", "TaxAmount",
		"TotalAmount", "PrintMark", "RandomNumber", "MainRemark", "CarrierType",
		"CarrierId1", "CarrierId2", "NPOBAN", "Description", "Quantity",
		"UnitPrice", "Amount", "DetailTaxType", "Remark",
	}

	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("writing header: %w", err)
	}

	// Count invoices and process rows
	invoiceCount := 0
	for _, row := range rows {
		if row.LineNo == 1 {
			invoiceCount++
		}

		// Check for rare Unicode characters in BuyerName
		buyerName := s.nullStringToString(row.BuyerName)
		if s.containsRareUnicode(buyerName) {
			log.Printf("BuyerName: %s contains rare characters", buyerName)
		}

		// Convert row to CSV record - handle NULLs properly
		record := []string{
			s.nullStringToString(row.InvoiceNumber),
			s.nullStringToString(row.InvoiceDate),
			s.nullStringToString(row.InvoiceTime),
			s.nullStringToString(row.BuyerIdentifier),
			buyerName,
			s.nullStringToString(row.BuyerAddress),
			s.nullStringToString(row.BuyerTelephoneNumber),
			s.nullStringToString(row.BuyerEmailAddress),
			s.nullStringToString(row.SalesAmount),
			s.nullStringToString(row.FreeTaxSalesAmount),
			s.nullStringToString(row.ZeroTaxSalesAmount),
			s.nullStringToString(row.TaxType),
			s.nullStringToString(row.TaxRate),
			s.nullStringToString(row.TaxAmount),
			s.nullStringToString(row.TotalAmount),
			s.nullStringToString(row.PrintMark),
			s.nullStringToString(row.RandomNumber),
			s.nullStringToString(row.MainRemark),
			s.nullStringToString(row.CarrierType),
			s.nullStringToString(row.CarrierID1),
			s.nullStringToString(row.CarrierID2),
			s.nullStringToString(row.NPOBAN),
			s.nullStringToString(row.Description),
			s.nullStringToString(row.Quantity),
			s.nullStringToString(row.UnitPrice),
			s.nullStringToString(row.Amount),
			s.nullStringToString(row.DetailTaxType),
			s.nullStringToString(row.Remark),
		}

		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("writing record: %w", err)
		}
	}

	// Write "Finish" marker
	finishRecord := []string{"Finish"}
	if err := writer.Write(finishRecord); err != nil {
		return nil, fmt.Errorf("writing finish marker: %w", err)
	}

	return &CSVGenerationResult{
		InvoiceCount: invoiceCount,
		TotalRows:    len(rows),
		FilePath:     csvFilePath,
	}, nil
}

// Helper function for null value handling
func (s *InvoiceService) nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// containsRareUnicode checks for rare Unicode characters (equivalent to your Python function)
func (s *InvoiceService) containsRareUnicode(text string) bool {
	for _, r := range text {
		if r >= 0x20000 && r <= 0x2FFFF {
			return true
		}
	}
	return false
}

// Gin handler example
func (s *InvoiceService) HandleGenC0401(c *gin.Context) {
	segmentNo := c.Param("segment_no")
	invoiceDate := c.Param("invoice_date")

	if segmentNo == "" || invoiceDate == "" {
		c.JSON(400, gin.H{"error": "segment_no and invoice_date are required"})
		return
	}

	result, err := s.GenC0401(segmentNo, invoiceDate)
	if err != nil {
		log.Printf("Error generating C0401: %v", err)
		c.JSON(500, gin.H{"error": "Failed to generate CSV"})
		return
	}

	c.JSON(200, gin.H{
		"invoice_count": result.InvoiceCount,
		"total_rows":    result.TotalRows,
		"file_path":     result.FilePath,
		"message":       "CSV generated successfully",
	})
}
