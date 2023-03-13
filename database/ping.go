package database

import (
	"go_redis/interface/resp"
	"go_redis/resp/reply"
)

func init() {
	RegisterCommand("ping", ping, 1)
}

func ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}
