package main

import (
	"fmt"
	"log/slog"
	"os"
)

type config struct {
	content []byte
}

// readConfig simulates reading a configuration file.
func readConfig(path string) (*config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed_to_read_config_file: %w", err)
	}
	return &config{content: content}, nil
}

// initApp initializes the application and returns an error if it fails.
func initApp() error {
	config, err := readConfig("config.txt")
	if err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}
	fmt.Println(string(config.content))
	return nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	err := initApp()
	if err != nil {
		logger.Error("App startup failed", "err", err)
		os.Exit(1)
	}
	logger.Info("App started successfully")
}
