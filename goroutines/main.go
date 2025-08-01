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
	mu     sync.Mutex // Mutex to protect the order status
}

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
)

var (
	totalUpdates int
	updateMutex  sync.Mutex
)
var OrdersStatuses = []string{StatusPending, StatusProcessing, StatusCompleted}

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

// func ProcessOrders(orders []*Order) {
// 	for _, order := range orders {
// 		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Simulate processing time
// 		fmt.Printf("Processing Order ID: %d\n", order.ID)
// 	}
// }

func UpdateOrders(order *Order) {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Simulate processing time
	order.Status = OrdersStatuses[rand.Intn(len(OrdersStatuses))]
	fmt.Printf("Updating Order ID: %d status to: %s\n", order.ID, order.Status)

	currentUpdates := totalUpdates
	time.Sleep(10 * time.Millisecond)
	totalUpdates = currentUpdates + 1
}

func ReportOrderStatuses(orders []*Order) {
	for _, order := range orders {
		fmt.Printf("OrderID: %d, Status: %s\n", order.ID, order.Status)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(3) // We will wait for two goroutines

	orders := GenerateOrders(10)

	// go func() {
	// 	defer wg.Done()
	// 	ProcessOrders(orders)
	// }()

	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			for _, order := range orders {
				UpdateOrders(order)
			}
		}()
	}

	wg.Wait()

	ReportOrderStatuses(orders)
	fmt.Printf("Total Updates %d", totalUpdates)
}
