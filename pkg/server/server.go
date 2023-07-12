package server

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	maxWorkers int
	powRepo    PoWRepo
	quoteRepo  QuoteRepo
	listener   net.Listener
	logger     *zap.Logger
}

func NewServer(
	maxWorkers int,
	powRepo PoWRepo,
	quoteRepo QuoteRepo,
	listener net.Listener,
	logger *zap.Logger,

) *Server {
	return &Server{
		maxWorkers: maxWorkers,
		powRepo:    powRepo,
		quoteRepo:  quoteRepo,
		listener:   listener,
		logger:     logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	go s.initiateShutdownOnSignal(shutdown, cancel)

	ch := make(chan net.Conn, s.maxWorkers)

	go func() {
		for {

			select {
			case <-ctx.Done():
				return
			case conn := <-ch:
				go s.handleConnection(ctx, conn, &wg)

			}
		}
	}()

	for {
		conn, err := s.acceptConnection()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				s.logger.Info("Listener has been closed successfully")
				break
			} else {
				s.logger.Error("Failed to accept new connection", zap.Error(err))
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(1 * time.Second):
					continue // continue accepting new connections even if an error occurs
				}
			}
		}

		ch <- conn
	}

	wg.Wait()
	<-ctx.Done()
	return ctx.Err()
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	wg.Add(1)

	defer func() {
		wg.Done()
		conn.Close()
	}()

	worker := NewWorker(conn, s.powRepo, s.quoteRepo, s.logger.With(zap.String("clientAddr", conn.RemoteAddr().String())))

	if err := worker.Handle(ctx); err != nil {
		s.logger.Error("Error handling connection", zap.Error(err))
	}
}

func (s *Server) initiateShutdownOnSignal(shutdown <-chan os.Signal, cancel context.CancelFunc) {
	<-shutdown
	s.logger.Info("Graceful shutdown initiated")
	s.listener.Close()
	cancel()
}

func (s *Server) acceptConnection() (net.Conn, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			s.logger.Info("Listener has been closed successfully")
		} else {
			s.logger.Error("Failed to accept new connection", zap.Error(err))
		}
	}
	return conn, err
}
