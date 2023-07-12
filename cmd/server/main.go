package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	"wisdom/pkg/pow"
	"wisdom/pkg/quote"
	"wisdom/pkg/server"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	log, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	err = godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	config, err := server.NewConfig()
	if err != nil {
		return fmt.Errorf("error creating config: %w", err)
	}

	powRepo := pow.NewPoW(config.Complexity)
	quoteRepo := quote.New(
		[]string{
			"Life is what happens when you're busy making other plans.",
			"Do not dwell in the past, do not dream of the future, concentrate the mind on the present moment.",
			"Life is really simple, but we insist on making it complicated.",
			"Love the life you live. Live the life you love.",
			"The only impossible journey is the one you never begin.",
			"Life is 10% what happens to you and 90% how you react to it.",
			"Life is like riding a bicycle. To keep your balance, you must keep moving.",
			"Life is a series of natural and spontaneous changes. Don't resist them - that only creates sorrow. Let reality be reality. Let things flow naturally forward in whatever way they like.",
		},
	)

	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		return fmt.Errorf("failed to start server listener: %w", err)
	}
	defer listener.Close()

	s := server.NewServer(
		config.MaxWorkers,
		powRepo,
		quoteRepo,
		listener,
		log,
	)

	err = s.Run(context.Background())
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Info("Server stopped by context cancellation")
		} else {
			return fmt.Errorf("server run error: %w", err)
		}
	}

	return nil
}
