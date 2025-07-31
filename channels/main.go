package main

import (
	"fmt"
)

func main() {
	ch2 := make(chan string, 2)

	go func() {
		ch2 <- "msg1"
		ch2 <- "msg2"
		// This will block if more than 2 messages are sent
		// until the main thread reads from the channel

		ch2 <- "msg3"
		fmt.Println("Buffered channel goroutine finished sending messages")
	}()

	fmt.Println(<-ch2, "and", <-ch2, "plus", <-ch2)

}
