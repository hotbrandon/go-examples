package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Message struct {
	// Chat ID can be acquired by calling curl https://api.telegram.org/bot<token>/getUpdates
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type TelegramBot struct {
	Token  string
	ChatID string
}

func NewTelegramBot(token, chatID string) *TelegramBot {
	return &TelegramBot{
		Token:  token,
		ChatID: chatID,
	}
}

func (t *TelegramBot) SendMessage(text string) error {
	requestURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Token)
	c := http.Client{Timeout: 10 * time.Second}

	// Prepare the message payload
	body, err := json.Marshal(Message{
		ChatID: t.ChatID,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Send the request
	resp, err := c.Post(requestURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil

}
