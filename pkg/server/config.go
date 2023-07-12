package server

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	MaxWorkers    int
	ServerAddress string
	Complexity    int
}

func NewConfig() (*Config, error) {
	maxWorkers, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))
	if err != nil {
		return nil, fmt.Errorf("error reading MAX_WORKERS from environment: %w", err)
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		return nil, errors.New("SERVER_ADDRESS must be set in environment")
	}

	complexity, err := strconv.Atoi(os.Getenv("COMPLEXITY"))
	if err != nil {
		return nil, fmt.Errorf("error reading COMPLEXITY from environment: %w", err)
	}

	return &Config{
		MaxWorkers:    maxWorkers,
		ServerAddress: serverAddress,
		Complexity:    complexity,
	}, nil
}
