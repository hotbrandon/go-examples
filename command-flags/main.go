package main

import (
	"flag"
	"fmt"
	"time"
)

// ./main -date=2025-10-31

type language = string

// From this point on, anywhere the compiler sees the identifier English, it
// substitutes it with the string value "en".
// Think of English as a readable, safe alias for "en".
const (
	English language = "en"
	French  language = "fr"
	Spanish language = "es"
)

func init() {
	fmt.Println("The init() function runs before main()")
}

func main() {
	today := time.Now().Format("2006-01-02")

	datePtr := flag.String("date", today, "invoice date in YYYY-MM-DD format")
	var userLanguage language
	flag.StringVar(&userLanguage, "lang", "en", "language")
	flag.Parse()

	fmt.Println("value of the -date flag:", *datePtr)

	switch userLanguage {
	case English:
		fmt.Println("Hello")
	case French:
		fmt.Println("Bonjour")
	case Spanish:
		fmt.Println("Hola")
	default:
		fmt.Println("Hello")
	}

}
