package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// StockData represents the structure of a single data object in the response.
// JSON tags are used to map the JSON keys to the struct fields.
type StockData struct {
	Date            string  `json:"date"`
	StockID         string  `json:"stock_id"`
	TradingVolume   int64   `json:"Trading_Volume"`
	TradingMoney    int64   `json:"Trading_money"`
	Open            float64 `json:"open"`
	Max             float64 `json:"max"`
	Min             float64 `json:"min"`
	Close           float64 `json:"close"`
	Spread          float64 `json:"spread"`
	TradingTurnover float32 `json:"Trading_turnover"`
}

// APIResponse represents the overall structure of the JSON response from the FinMind API.
type APIResponse struct {
	Msg    string      `json:"msg"`
	Status int         `json:"status"`
	Data   []StockData `json:"data"`
}

func main() {
	c := http.Client{
		Timeout: 10 * time.Second,
	}

	// The endpoint includes the path /api/v4/data
	req, err := http.NewRequest("GET", "https://api.finmindtrade.com/api/v4/data", nil)
	if err != nil {
		fmt.Printf("error creating request: %s\n", err)
		return
	}

	// Set up query parameters
	q := req.URL.Query()
	q.Add("dataset", "TaiwanStockPrice")
	q.Add("data_id", "2337")
	q.Add("start_date", "2025-10-31") // Example start date
	q.Add("end_date", "2025-10-31")   // Example end date

	// Encode the parameters and attach them to the request URL
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.RawQuery)

	fmt.Println("Sending request to:", req.URL.String())

	// Send the request
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error making request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body: %s\n", err)
		return
	}

	// Create a variable to hold our parsed data.
	var apiResponse APIResponse

	// Unmarshal the JSON body into our struct.
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Printf("error unmarshalling JSON: %s\n", err)
		return
	}

	// Now you can access the data in a structured way.
	fmt.Printf("Parsed Response: %+v\n", apiResponse)
}
