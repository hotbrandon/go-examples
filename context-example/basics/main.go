package main

import (
	"context"
	"time"
)

func doSleep(ctx context.Context, cancel context.CancelFunc) {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	doSleep(ctx, cancel)
}
