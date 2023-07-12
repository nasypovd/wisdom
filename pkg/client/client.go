package client

import (
	"encoding/json"
	"fmt"
	"net"
	"wisdom/pkg/domain"
	"wisdom/pkg/dto"

	"go.uber.org/zap"
)

type Solver interface {
	Solve(challenge domain.Challenge) domain.Solution
}

type Transport interface {
	ReceiveChallenge() (domain.Challenge, error)
	SendSolution(solution domain.Solution) error
	ReceiveQuote() (domain.Quote, error)
}

type Client struct {
	conn   net.Conn
	solver Solver
	logger *zap.Logger
}

func NewClient(conn net.Conn, solver Solver, log *zap.Logger) *Client {

	return &Client{
		conn:   conn,
		solver: solver,
		logger: log,
	}
}

func (c *Client) Run() (domain.Quote, error) {

	challenge, err := c.ReceiveChallenge()
	if err != nil {
		return "", fmt.Errorf("reading challenge from server failed: %w", err)
	}

	c.logger.Info("Received challenge", zap.String("challenge", challenge.Value), zap.Int("difficulty", challenge.Difficulty))

	solution := c.solver.Solve(challenge)

	c.SendSolution(solution)
	if err != nil {
		return "", fmt.Errorf("sending solution to server failed: %w", err)
	}

	quote, err := c.ReceiveQuote()
	if err != nil {
		return "", fmt.Errorf("reading quote from server failed: %w", err)
	}
	c.logger.Info("Received quote", zap.String("quote", string(quote)))

	return quote, nil
}

func (c *Client) ReceiveChallenge() (domain.Challenge, error) {
	var msg dto.Challenge
	dec := json.NewDecoder(c.conn)
	err := dec.Decode(&msg)
	if err != nil {
		return domain.Challenge{}, fmt.Errorf("failed to receive challenge: %w", err)
	}
	return domain.Challenge{
		Value:      msg.Value,
		Difficulty: msg.Difficulty,
	}, nil
}

func (c *Client) SendSolution(solution domain.Solution) error {
	msg := transport.Solution{
		Body: string(solution),
	}

	enc := json.NewEncoder(c.conn)
	err := enc.Encode(msg)
	if err != nil {
		return fmt.Errorf("failed to send solution: %w", err)
	}

	return nil
}

func (c *Client) ReceiveQuote() (domain.Quote, error) {
	var msg transport.Quote
	dec := json.NewDecoder(c.conn)
	err := dec.Decode(&msg)
	if err != nil {
		return "", fmt.Errorf("failed to receive quote: %w", err)
	}

	return domain.Quote(msg.Body), nil
}
