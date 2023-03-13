package connection

import (
	"go_redis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type Connection struct {
	conn         net.Conn
	waitingReply wait.Wait
	mu           sync.Mutex
	selectDB     int
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) Write(bytes []byte) (err error) {
	if len(bytes) == 0 {
		return nil
	}
	defer func() {
		c.mu.Unlock()
		c.waitingReply.Done()
	}()
	c.mu.Lock()
	c.waitingReply.Add(1)
	_, err = c.conn.Write(bytes)
	return
}
func (c *Connection) GetDBIndex() int {
	return c.selectDB
}

func (c *Connection) SelectDB(dbIndex int) {
	c.selectDB = dbIndex
}

func (c *Connection) Close() (err error) {
	c.waitingReply.WaitWithTimeout(5 * time.Second)
	c.conn.Close()
	return
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.RemoteAddr()
}
