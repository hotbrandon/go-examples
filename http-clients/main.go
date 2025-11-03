package main

import (
	"fmt"
	"log"
)

func main() {
	//begin_date := time.Now().Format("2006-01-02")
	begin_date := "2025-11-04"
	ticker := "2337"
	dataset := "TaiwanStockPrice"
	f := NewFinMind(dataset, ticker, begin_date, begin_date)

	data, err := f.GetTaiwanStockPrice()

	// Unmarshal the JSON body into our struct.
	if err != nil {
		// If we only care about success/failure, we can use a simple message.
		log.Fatal("Failed to get stock quote.")
	}

	// Now you can access the data in a structured way.
	fmt.Printf("Close Price for %s on %s is %.2f\n", ticker, data.Date, data.Close)
}
