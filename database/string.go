package database

import (
	"go_redis/interface/database"
	"go_redis/interface/resp"
	"go_redis/lib/utils"
	"go_redis/resp/reply"
)

// GET
func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	if val, ok := db.GetEntity(key); !ok {
		return reply.MakeNullBulkReply()
	} else {
		entity := val.Data.([]byte)
		return reply.MakeBulkReply(entity)
	}
}

// SET
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	dataEntity := &database.DataEntity{
		Data: val,
	}
	db.PutEntity(key, dataEntity)
	db.addAof(utils.ToCmdLine2("set", args...))
	return reply.MakeOkReply()
}

// SETNX
func execSetNx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	dataEntity := &database.DataEntity{
		Data: val,
	}
	result := db.PutIfAbsent(key, dataEntity)
	db.addAof(utils.ToCmdLine2("setnx", args...))
	return reply.MakeIntReply(int64(result))
}

// GETSET
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity, ok := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{Data: val})
	if !ok {
		return reply.MakeNullBulkReply()
	} else {
		entityByte := entity.Data.([]byte)
		db.addAof(utils.ToCmdLine2("getset", args...))
		return reply.MakeBulkReply(entityByte)
	}
}

// STRLEN
func execStrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	if entity, ok := db.GetEntity(key); !ok {
		return reply.MakeNullBulkReply()
	} else {
		bytes := entity.Data.([]byte)
		return reply.MakeIntReply(int64(len(bytes)))
	}
}

func init() {
	RegisterCommand("GET", execGet, 2)
	RegisterCommand("SET", execSet, 3)
	RegisterCommand("SETNX", execSetNx, 3)
	RegisterCommand("GETSET", execGetSet, 3)
	RegisterCommand("STRLEN", execStrLen, 2)
}
