package resphandler

import (
	"context"
	"go-redis/db/database"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"
	"go-redis/utils/sync/atomic"
	"go.uber.org/zap"
	"io"
	"net"
	"strings"
	"sync"
)

// ================================================================================

type Config struct {
	Database database.Interface
}

type Handler struct {
	activeConn sync.Map
	closing    atomic.Bool
	db         database.Interface
}

func NewHandlerWithDB(config *Config) *Handler {
	return &Handler{db: config.Database}
}

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		_ = conn.Close()
	}

	client := connection.New(conn)
	h.activeConn.Store(client, struct{}{})

	// handle tcp
	ch := parser.ParseStream(conn)
	for stream := range ch {
		// error
		if stream.Err != nil {
			if stream.Err == io.EOF || stream.Err == io.ErrUnexpectedEOF ||
				strings.Contains(stream.Err.Error(), "use of closed network connection") {
				h.closeClient(client)
				zap.S().Infof("connection closed: %s", client.RemoteAddr().String())
				return
			}
			// protocol error
			errReply := reply.NewErrReply(stream.Err.Error())
			if err := client.Write(errReply.ToBytes()); err != nil {
				h.closeClient(client)
				zap.S().Infof("connection closed: %s", client.RemoteAddr().String())
				return
			}
			continue
		}
		// no error
		if stream.Data == nil {
			continue
		}

		multiBulkReply, ok := stream.Data.(*reply.MultiBulkReply)
		if !ok {
			zap.S().Error("require multi bulk reply")
			continue
		}
		execReply := h.db.Exec(client, multiBulkReply.Args)
		if execReply != nil {
			_ = client.Write(execReply.ToBytes())
		} else {
			_ = client.Write(reply.NewUnknowErrReply().ToBytes())
		}
	}
}

// close one client
func (h *Handler) closeClient(client *connection.Connection) {
	_ = client.Close()
	h.db.AfterClientClose(client)
	h.activeConn.Delete(client)
}

// Close handler and all client
func (h *Handler) Close() error {
	zap.S().Info("handler shutting down...")
	h.closing.Set(true)
	h.activeConn.Range(func(key, value any) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	h.db.Close()
	return nil
}
