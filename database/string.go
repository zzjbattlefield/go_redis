package database

import (
	"go_redis/interface/database"
	"go_redis/interface/resp"
	"go_redis/lib/utils"
	"go_redis/resp/reply"
	"strconv"
	"strings"
	"time"
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

const (
	upsertPolicy = iota //default set
	insertPolict        // set nx
	updatePolicy        // set xx
)

const unlimitTTL int64 = 0

// SET
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	policy := upsertPolicy
	ttl := unlimitTTL
	if len(args) > 2 {
		for i := 2; i < len(args); i++ {
			arg := strings.ToUpper(string(args[i]))
			if arg == "EX" {
				if ttl != unlimitTTL || i+1 > len(args) {
					return &reply.SyntaxErrReply{}
				}
				ttlArg, err := strconv.ParseInt(string(args[i+1]), 10, 64)
				if err != nil {
					return &reply.SyntaxErrReply{}
				}
				ttl = ttlArg * 1000 //转成毫秒
				i++
			} else if arg == "PX" {
				if ttl != unlimitTTL || i+1 > len(args) {
					return &reply.SyntaxErrReply{}
				}
				ttlArg, err := strconv.ParseInt(string(args[i+1]), 10, 64)
				if err != nil {
					return &reply.SyntaxErrReply{}
				}
				ttl = ttlArg
				i++
			} else if arg == "NX" {
				if policy == updatePolicy {
					return &reply.SyntaxErrReply{}
				}
				policy = insertPolict
			} else if arg == "XX" {
				if policy == insertPolict {
					return &reply.SyntaxErrReply{}
				}
				policy = updatePolicy
			}
		}
	}
	dataEntity := &database.DataEntity{
		Data: val,
	}
	var result int
	switch policy {
	case upsertPolicy:
		result = db.PutEntity(key, dataEntity)
	case insertPolict:
		result = db.PutIfAbsent(key, dataEntity)
	case updatePolicy:
		result = db.PutIfExists(key, dataEntity)
	}
	if result > 0 {
		if ttl != unlimitTTL {
			//设置过期时间
			db.Expire(key, time.Now().Add(time.Duration(ttl)*time.Millisecond))
		}
		return reply.MakeOkReply()
	}
	return reply.MakeNullBulkReply()
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
	RegisterCommand("SET", execSet, -3)
	RegisterCommand("SETNX", execSetNx, 3)
	RegisterCommand("GETSET", execGetSet, 3)
	RegisterCommand("STRLEN", execStrLen, 2)
}
