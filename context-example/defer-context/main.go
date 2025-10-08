package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("start:", runtime.NumGoroutine())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done() // wait until canceled or timeout
	}()

	time.Sleep(1 * time.Second)
	fmt.Println("before cancel (still blocked):", runtime.NumGoroutine()) // typically 2

	cancel()  // cancel early (don’t wait for the 5s timeout)
	wg.Wait() // ensure the goroutine exits
	fmt.Println("after wg.Wait():", runtime.NumGoroutine())
	// give the scheduler a moment (usually not necessary after wg.Wait(), but harmless)
	time.Sleep(10 * time.Millisecond)

	fmt.Println("finally:", runtime.NumGoroutine()) // typically back to 1
}
