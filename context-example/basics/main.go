// ...existing code...
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
)

type userKey string

func main() {
	// parse a single numeric arg: 1=Deadline, 2=Cancellation, 3=Values
	if len(os.Args) < 2 {
		usage()
		return
	}
	mode := os.Args[1]
	// accept numeric or text form
	switch mode {
	case "1", "deadline":
		exampleDeadline()
	case "2", "cancel":
		exampleCancel()
	case "3", "values":
		exampleValues()
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage: go run . <mode>")
	fmt.Println("Modes:")
	fmt.Println("  1 | deadline  - demonstrate context deadline/timeout")
	fmt.Println("  2 | cancel    - demonstrate cancellation signals")
	fmt.Println("  3 | values    - demonstrate request-scoped values")
	fmt.Println("Example (Windows PowerShell/CMD):")
	fmt.Println("  go run . 1")
}

// Deadline example: deadline shorter than work -> context times out
func exampleDeadline() {
	fmt.Println("Example 1: Deadline/Timeout")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case <-time.After(2 * time.Second): // simulated work longer than deadline
		fmt.Println("work finished (unexpected)")
	case <-ctx.Done():
		fmt.Println("context done:", ctx.Err()) // expected: context deadline exceeded
	}
}

// Cancellation signal example: start work, cancel from parent
func exampleCancel() {
	fmt.Println("Example 2: Cancellation Signal")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // good practice even if cancel is called below

	// cancel after 1s from a separate goroutine
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("parent: calling cancel()")
		cancel()
	}()

	// worker watches ctx.Done()
	for i := 1; i <= 5; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("worker: stopped early:", ctx.Err())
			return
		case <-time.After(500 * time.Millisecond):
			fmt.Println("worker: did unit", i)
		}
	}
	fmt.Println("worker: completed all units")
}

// Request-scoped values example: pass values via context
func exampleValues() {
	fmt.Println("Example 3: Request-scoped values")

	// attach values to context
	ctx := context.WithValue(context.Background(), userKey("requestID"), "req-42")
	ctx = context.WithValue(ctx, userKey("user"), "alice")
	handleRequest(ctx)
}

func handleRequest(ctx context.Context) {
	// retrieve values
	reqID, _ := ctx.Value(userKey("requestID")).(string)
	user, _ := ctx.Value(userKey("user")).(string)

	fmt.Println("handleRequest: requestID =", reqID)
	fmt.Println("handleRequest: user =", user)

	// show that cancelling or timing out would be picked up if present
	// (here we just demonstrate reading values)
	_ = strconv.Itoa // keep import use explicit if needed elsewhere
}

// ...existing code...
