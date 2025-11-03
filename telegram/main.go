package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN environment variable is not set")
	}

	chatID := os.Getenv("CHAT_ID")
	if chatID == "" {
		log.Fatal("CHAT_ID environment variable is not set")
	}
	tgBot := NewTelegramBot(token, chatID)

	err := tgBot.SendMessage("Go Fish!!")
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}

	fmt.Println("message sent successfully!")

}
