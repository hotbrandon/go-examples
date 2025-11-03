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

type FinMind struct {
	dataset    string
	data_id    string
	start_date string
	end_date   string
}

func NewFinMind(dataset, data_id, start_date, end_date string) *FinMind {
	return &FinMind{
		dataset:    dataset,
		data_id:    data_id,
		start_date: start_date,
		end_date:   end_date,
	}
}

func (f *FinMind) GetTaiwanStockPrice() (*StockData, error) {
	c := http.Client{
		Timeout: 10 * time.Second,
	}

	// The endpoint includes the path /api/v4/data
	req, err := http.NewRequest("GET", "https://api.finmindtrade.com/api/v4/data", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)

	}

	// Set up query parameters
	q := req.URL.Query()
	q.Set("dataset", f.dataset)
	q.Set("data_id", f.data_id)
	q.Set("start_date", f.start_date)
	q.Set("end_date", f.end_date)

	// Encode the parameters and attach them to the request URL
	req.URL.RawQuery = q.Encode()

	fmt.Println("Sending request to:", req.URL.String())

	// Send the request
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)

	}

	// Create a variable to hold our parsed data.
	var apiResponse APIResponse

	// Unmarshal the JSON body into our struct.
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)

	}

	if len(apiResponse.Data) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	return &apiResponse.Data[0], nil
}
