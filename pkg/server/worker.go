package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"
	"wisdom/pkg/domain"
	"wisdom/pkg/dto"

	"go.uber.org/zap"
)

const MAX_RETRY = 3
const RETRY_SLEEP = 2 * time.Second

type PoWRepo interface {
	Generate() domain.Challenge
	Verify(challenge domain.Challenge, solution domain.Solution) bool
}

type QuoteRepo interface {
	Get() domain.Quote
}

type Worker struct {
	conn      net.Conn
	pow       PoWRepo
	quoteRepo QuoteRepo
	logger    *zap.Logger
}

func NewWorker(conn net.Conn, pow PoWRepo, quoteRepo QuoteRepo, logger *zap.Logger) *Worker {
	return &Worker{conn: conn, pow: pow, quoteRepo: quoteRepo, logger: logger}
}

func (w *Worker) Handle(ctx context.Context) error {
	challenge := w.pow.Generate()
	w.logger.
		With(zap.String("challenge", challenge.Value)).
		With(zap.Int("difficulty", challenge.Difficulty)).
		Info("Generated challenge")

	err := w.retry(ctx, func() error {
		return w.SendChallenge(challenge)
	}, "sending challenge to client")
	if err != nil {
		return err
	}

	var solution domain.Solution
	err = w.retry(ctx, func() error {
		var err error
		solution, err = w.ReceiveSolution()
		return err
	}, "reading solution from client")
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if w.pow.Verify(challenge, solution) {
			w.logger.Info("Challenge solved. Sending quote.")
			quote := w.quoteRepo.Get()
			err := w.SendQuote(quote)
			if err != nil {
				return fmt.Errorf("sending quote to client failed: %w", err)
			}
		} else {
			w.logger.Error("Challenge not solved. Closing connection.")
		}
	}

	return nil
}

func (w *Worker) retry(ctx context.Context, action func() error, actionDescription string) error {
	for i := 1; i <= MAX_RETRY; i++ {
		err := action()
		if err == nil {
			return nil
		}
		if i == MAX_RETRY {
			return fmt.Errorf("%s failed: %w", actionDescription, err)
		}
		w.logger.Warn(fmt.Sprintf("Failed to %s, retrying...", actionDescription), zap.Error(err))

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(RETRY_SLEEP):
		}
	}
	return nil
}

func (w *Worker) SendChallenge(challenge domain.Challenge) error {
	msg := dto.Challenge{
		Value:      challenge.Value,
		Difficulty: challenge.Difficulty,
	}

	enc := json.NewEncoder(w.conn)
	err := enc.Encode(msg)
	if err != nil {
		return fmt.Errorf("failed to send challenge: %w", err)
	}

	return nil
}

func (w *Worker) ReceiveSolution() (domain.Solution, error) {
	var msg dto.Solution
	dec := json.NewDecoder(w.conn)
	err := dec.Decode(&msg)
	if err != nil {
		return "", fmt.Errorf("failed to receive solution: %w", err)
	}

	return domain.Solution(msg.Body), nil
}

func (w *Worker) SendQuote(quote domain.Quote) error {
	msg := dto.Quote{
		Body: string(quote),
	}

	enc := json.NewEncoder(w.conn)
	err := enc.Encode(msg)
	if err != nil {
		return fmt.Errorf("failed to send quote: %w", err)
	}

	return nil
}
