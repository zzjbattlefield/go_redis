package handler

import (
	"context"
	"go_redis/database"
	databaseface "go_redis/interface/database"
	"go_redis/lib/logger"
	"go_redis/lib/sync/atomic"
	"go_redis/resp/connection"
	"go_redis/resp/parser"
	"go_redis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
)

type RespHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
	db         databaseface.Database
}

func (r *RespHandler) closeClient(client *connection.Connection) {
	client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

func NewHandler() *RespHandler {
	return &RespHandler{
		db: database.NewDatabase(),
	}
}

func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		conn.Close()
	}
	client := connection.NewConnection(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parser.ParserStream(conn)
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF || payload.Err == io.ErrUnexpectedEOF || strings.Contains(payload.Err.Error(), "use of closed network connection") {
				//tcp挥手 或者 使用关闭的连接
				r.closeClient(client)
				logger.Info("connect closed:" + payload.Err.Error())
				return
			}
			//协议错误
			errReply := reply.MakeErrReply(payload.Err.Error())
			if err := client.Write(errReply.ToByte()); err != nil {
				r.closeClient(client)
				logger.Info("connect closed:" + client.RemoteAddr().String())
				return
			}
			continue
		}
		//exec
		if payload.Data == nil {
			continue
		}
		replyData, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		if result := r.db.Exec(client, replyData.Args); result != nil {
			client.Write(result.ToByte())
		} else {
			client.Write(reply.MakeErrReply("ERR unknow").ToByte())
		}

	}
}
func (r *RespHandler) Close() (err error) {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	r.activeConn.Range(func(key, value any) bool {
		client := key.(*connection.Connection)
		client.Close()
		return true
	})
	r.db.Close()
	return
}
