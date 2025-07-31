package main

import (
	"fmt"
	"time"
)

var counter int

func increment() {
	for i := 0; i < 1000; i++ {
		time.Sleep(1 * time.Millisecond)
		counter++
	}
}

func main() {
	for i := 0; i < 5; i++ {
		go increment() // Start 5 goroutines to increment the counter
	}

	time.Sleep(10 * time.Second)
	fmt.Println("Final counter value:", counter)
}
