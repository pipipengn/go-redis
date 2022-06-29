package tcp

import (
	"bufio"
	"context"
	"go-redis/utils/sync/atomic"
	"go-redis/utils/sync/wait"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"time"
)

// Client one client
type Client struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (c *Client) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	_ = c.Conn.Close()
	return nil
}

// Handler tcp server handler impl
type Handler struct {
	activeConn sync.Map
	closing    atomic.Bool
}

func NewHandler() *Handler {
	return &Handler{}
}

// Handle concurrently handle each conn
func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		_ = conn.Close()
	}

	// store all connected client
	client := &Client{Conn: conn}
	h.activeConn.Store(client, struct{}{})

	// handle tcp
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				zap.S().Info("Connection closed")
				h.activeConn.Delete(client)
			} else {
				zap.S().Warn(err)
			}
			return
		}

		client.Waiting.Add(1)
		_, _ = conn.Write([]byte(msg))
		client.Waiting.Done()
	}
}

func (h *Handler) Close() error {
	zap.S().Info("handler is shutting down...")
	h.closing.Set(true)
	h.activeConn.Range(func(key, value any) bool {
		client := key.(*Client)
		_ = client.Close()
		return true
	})
	return nil
}
