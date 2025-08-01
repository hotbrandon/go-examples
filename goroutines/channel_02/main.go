package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	ID     int
	Status string
}

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
)

func GenerateOrders(n int) []*Order {
	var orders []*Order = make([]*Order, n)

	for i := 0; i < n; i++ {
		orders[i] = &Order{
			ID:     i + 1,
			Status: StatusPending,
		}
	}

	return orders
}

func ProcessOrders(orders []*Order) {
	for _, order := range orders {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Simulate processing time
		fmt.Printf("Processing Order ID: %d\n", order.ID)
	}
}

func ReportOrderStatuses(orders []*Order) {
	for _, order := range orders {
		fmt.Printf("OrderID: %d, Status: %s\n", order.ID, order.Status)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	orderChan := make(chan *Order, 10) // Buffered channel to hold orders

	go func() {
		defer wg.Done()
		for _, order := range GenerateOrders(20) {
			orderChan <- order
		}
		close(orderChan)
		fmt.Println("All orders generated and sent to channel.")
	}()

	go func() {
		defer wg.Done()
		for order := range orderChan {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond) // Simulate processing time
			order.Status = StatusProcessing                             // Simulate status update
			fmt.Printf("Processing Order ID: %d, Status: %s\n", order.ID, order.Status)
		}
	}()

	wg.Wait()
}
