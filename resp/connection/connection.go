package connection

import (
	"go-redis/utils/sync/wait"
	"net"
	"sync"
	"time"
)

type Connection struct {
	conn       net.Conn
	waiting    wait.Wait
	mu         sync.Mutex
	selectedDB int
}

func New(conn net.Conn) *Connection {
	return &Connection{conn: conn}
}

func NewEmpty() *Connection {
	return &Connection{}
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	c.mu.Lock()
	c.waiting.Add(1)
	defer func() {
		c.waiting.Done()
		c.mu.Unlock()
	}()

	if _, err := c.conn.Write(bytes); err != nil {
		return err
	}
	return nil
}

func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}

func (c *Connection) SelectDB(i int) {
	c.selectedDB = i
}

func (c *Connection) Close() error {
	c.waiting.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
	return nil
}
