package main

import (
	"flag"
	"fmt"
	"time"
)

// ./main -date=2025-10-31

func init() {
	fmt.Println("The init() function runs before main()")
}

func main() {
	today := time.Now().Format("2006-01-02")

	datePtr := flag.String("date", today, "invoice date in YYYY-MM-DD format")
	flag.Parse()

	fmt.Println("value of the -date flag:", *datePtr)

}
