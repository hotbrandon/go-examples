package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	botToken := "mytoken"  // replace with your bot token
	chatID := "8339256565" // replace with your chat ID
	message := "awesome!"

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	// Telegram API expects JSON payload
	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)
}
