package database

import (
	"go_redis/interface/resp"
	"go_redis/resp/reply"
)

type EchoDatabase struct {
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

func (e *EchoDatabase) Exec(client resp.Connect, data resp.Reply) resp.Reply {
	switch data.(type) {
	case *reply.MultiBulkReply:
		args := data.(*reply.MultiBulkReply).Args
		return reply.MakeMultiBulkReply(args)
	case *reply.BulkReply:
		args := data.(*reply.BulkReply).Args
		return reply.MakeBulkReply(args)
	case *reply.IntReply:
		code := data.(*reply.IntReply).Code
		return reply.MakeIntReply(code)
	case *reply.StatusReply:
		status := data.(*reply.StatusReply).Status
		return reply.MakeStatusReply(status)
	}
	return &reply.NoReply{}
}
func (e *EchoDatabase) Close() {

}
func (e *EchoDatabase) AfterClientClose(c resp.Connect) {

}
