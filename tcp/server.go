package tcp

import (
	"context"
	"go.uber.org/zap"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type HandlerInterface interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

// Server go-redis server
type Server struct {
	address string
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
	}
}

func (s *Server) ListenAndServeWithSignal(handler HandlerInterface) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		zap.S().Error("cannot listen tcp: %v", zap.Error(err))
		return err
	}
	zap.S().Info("start tcp listen")

	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	if err = s.listenAndServe(listener, handler, closeChan); err != nil {
		zap.S().Error("cannot serve tcp: %v", zap.Error(err))
		return err
	}

	return nil
}

func (s *Server) listenAndServe(listener net.Listener, handler HandlerInterface, closeChan <-chan struct{}) error {
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	go func() {
		<-closeChan
		zap.S().Info("closing...")
		_ = listener.Close()
		_ = handler.Close()
	}()

	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		zap.S().Info("accepted link")
		wg.Add(1)

		go func() {
			defer wg.Done()
			handler.Handle(context.Background(), conn)
		}()
	}

	wg.Wait()
	return nil
}
