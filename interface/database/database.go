package database

import "go_redis/interface/resp"

type CmdLine = [][]byte

type Database interface {
	Exec(client resp.Connect, args [][]byte) resp.Reply
	Close()
	AfterClientClose(c resp.Connect)
}

type DataEntity struct {
	Data interface{}
}
