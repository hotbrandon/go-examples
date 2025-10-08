package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AdaptiveCard represents a basic adaptive card structure
type AdaptiveCard struct {
	Type    string        `json:"type"`
	Version string        `json:"version"`
	Body    []CardElement `json:"body"`
}

type CardElement struct {
	Type   string `json:"type"`
	Text   string `json:"text"`
	Weight string `json:"weight,omitempty"`
	Size   string `json:"size,omitempty"`
}

type AdaptiveCardMessage struct {
	Type        string       `json:"type"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	ContentType string       `json:"contentType"`
	Content     AdaptiveCard `json:"content"`
}

func sendAdaptiveCard(webhookURL, title, text string) error {
	card := AdaptiveCard{
		Type:    "AdaptiveCard",
		Version: "1.2",
		Body: []CardElement{
			{
				Type:   "TextBlock",
				Text:   title,
				Weight: "Bolder",
				Size:   "Medium",
			},
			{
				Type: "TextBlock",
				Text: text,
			},
		},
	}

	payload := AdaptiveCardMessage{
		Type: "message",
		Attachments: []Attachment{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				Content:     card,
			},
		},
	}

	return sendMessage(webhookURL, payload)
}

func sendMessage(webhookURL string, payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Power Automate returns 202 (Accepted) for successful requests
	if resp.StatusCode != http.StatusOK && resp.StatusCode != 202 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	// Replace with your actual Power Automate webhook URL
	webhookURL := "https://defaulte7c12ff53b3c436da10ac8ce9f3c00.a7.environment.api.powerplatform.com:443/powerautomate/automations/direct/workflows/25c8ef20e0434cadb6ab3c00053fc555/triggers/manual/paths/invoke?api-version=1&sp=%2Ftriggers%2Fmanual%2Frun&sv=1.0&sig=LLVe4gRCGlmpfLE2MaaorDWTwoSRnIuOt3K2b-S8Aow"

	fmt.Println("Testing Teams notifications via Power Automate...")

	fmt.Println("Sending detailed Adaptive Card...")
	if err := sendAdaptiveCard(webhookURL, "Detailed Report", "This is a more detailed adaptive card with multiple text blocks."); err != nil {
		fmt.Printf("Error sending adaptive card: %v\n", err)
	} else {
		fmt.Println("Detailed adaptive card sent successfully!")
	}
}
