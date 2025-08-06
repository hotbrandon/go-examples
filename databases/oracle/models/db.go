package models

import (
	"database/sql"
	"log"
	"strings"
)

type OracleConfig struct {
	DSN string
}

var OracleConfigs = make(map[string]OracleConfig)

func GetOracleConnection(segment_no string) (*sql.DB, error) {
	// the config will always be available because it is set in main
	config := OracleConfigs[strings.ToUpper(segment_no)]
	db, err := sql.Open("oracle", config.DSN)
	if err != nil {
		log.Printf("Error opening Oracle connection: %v", err)

		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging Oracle database: %v", err)
		return nil, err
	}
	return db, nil
}

type Invoice struct {
	InvoiceID string
}

// C0401Data represents a single row of data from the MRIF004 query,
// which can be either a main invoice record or a detail line.
// The struct fields are ordered to match the C0401 CSV format.
type C0401Data struct {
	// --- Invoice Header Fields ---
	// These fields are populated by the database query only for the first line of an invoice.
	// For subsequent lines, they will be nil.

	// InvoiceNumber is the invoice identifier. It's always present.
	InvoiceNumber string `db:"invoicenumber" csv:"InvoiceNumber"`

	// InvoiceDate is the date of the invoice (YYYYMMDD). Pointer to handle NULLs.
	InvoiceDate *string `db:"invoicedate" csv:"InvoiceDate"`

	// InvoiceTime is the time of the invoice. Pointer to handle NULLs.
	InvoiceTime *string `db:"invoicetime" csv:"InvoiceTime"`

	// BuyerIdentifier is the buyer's tax ID. Pointer to handle NULLs.
	BuyerIdentifier *string `db:"buyeridentifier" csv:"BuyerIdentifier"`

	// BuyerName is the name of the buyer. Pointer to handle NULLs.
	BuyerName *string `db:"buyername" csv:"BuyerName"`

	// BuyerAddress is the address of the buyer. Pointer to handle NULLs.
	BuyerAddress *string `db:"buyeraddress" csv:"BuyerAddress"`

	// BuyerTelephoneNumber is the buyer's phone number. Pointer to handle NULLs.
	BuyerTelephoneNumber *string `db:"buyertelephonenumber" csv:"BuyerTelephoneNumber"`

	// BuyerEmailAddress is the buyer's email. Pointer to handle NULLs.
	BuyerEmailAddress *string `db:"buyeremailaddress" csv:"BuyerEmailAddress"`

	// SalesAmount is the total sales amount. Pointer to handle NULLs.
	SalesAmount *float64 `db:"salesamount" csv:"SalesAmount,omitempty"`

	// FreeTaxSalesAmount is the tax-free sales amount. Pointer to handle NULLs.
	FreeTaxSalesAmount *float64 `db:"freetaxsalesamount" csv:"FreeTaxSalesAmount,omitempty"`

	// ZeroTaxSalesAmount is the zero-tax sales amount. Pointer to handle NULLs.
	ZeroTaxSalesAmount *float64 `db:"zerotaxsalesamount" csv:"ZeroTaxSalesAmount,omitempty"`

	// TaxType indicates the tax type (e.g., '1' for taxable). Pointer to handle NULLs.
	TaxType *string `db:"taxtype" csv:"TaxType"`

	// TaxRate is fetched as a formatted string from the DB to preserve formatting (e.g., '0.99').
	// The Python code also treats it as a string. Pointer to handle NULLs.
	TaxRate *string `db:"taxrate" csv:"TaxRate"`

	// TaxAmount is the total tax amount. Pointer to handle NULLs.
	TaxAmount *float64 `db:"taxamount" csv:"TaxAmount,omitempty"`

	// TotalAmount is the grand total (SalesAmount + TaxAmount). Pointer to handle NULLs.
	TotalAmount *float64 `db:"totalamount" csv:"TotalAmount,omitempty"`

	// PrintMark indicates if the invoice is marked for printing. Pointer to handle NULLs.
	PrintMark *string `db:"printmark" csv:"PrintMark"`

	// RandomNumber is the invoice's random number for verification. Pointer to handle NULLs.
	RandomNumber *string `db:"randomnumber" csv:"RandomNumber"`

	// MainRemark contains any primary remarks for the invoice. Pointer to handle NULLs.
	MainRemark *string `db:"mainremark" csv:"MainRemark"`

	// CarrierType is the type of electronic carrier used. Pointer to handle NULLs.
	CarrierType *string `db:"carriertype" csv:"CarrierType"`

	// CarrierId1 is the primary carrier ID. Pointer to handle NULLs.
	CarrierId1 *string `db:"carrierid1" csv:"CarrierId1"`

	// CarrierId2 is the secondary carrier ID. Pointer to handle NULLs.
	CarrierId2 *string `db:"carrierid2" csv:"CarrierId2"`

	// NPOBAN is the donation code for non-profit organizations. Pointer to handle NULLs.
	NPOBAN *string `db:"npoban" csv:"NPOBAN"`

	// --- Invoice Detail Fields ---
	// These fields are populated for every line of an invoice.

	// Description of the item or service. Pointer to handle NULLs.
	Description *string `db:"description" csv:"Description"`

	// Quantity of the item. Assumed to be non-nullable based on Python code.
	Quantity float64 `db:"quantity" csv:"Quantity"`

	// UnitPrice of the item. Assumed to be non-nullable.
	UnitPrice float64 `db:"unitprice" csv:"UnitPrice"`

	// Amount is the total for the line (Quantity * UnitPrice). Assumed to be non-nullable.
	Amount float64 `db:"amount" csv:"Amount"`

	// DetailTaxType is the tax type for the specific detail line. Pointer to handle NULLs.
	DetailTaxType *string `db:"detailtaxtype" csv:"DetailTaxType"`

	// Remark contains any remarks for the detail line. Pointer to handle NULLs.
	Remark *string `db:"remark" csv:"Remark"`

	// --- Internal Control Field ---
	// LineNo is used for logic (identifying the header row) but not included in the CSV output.
	LineNo int `db:"line_no" csv:"-"`
}
