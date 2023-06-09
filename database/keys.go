package database

import (
	"go_redis/aof"
	"go_redis/interface/resp"
	"go_redis/lib/utils"
	"go_redis/resp/reply"
	"strconv"
	"time"
)

func init() {
	RegisterCommand("del", execDel, -2)
	RegisterCommand("exists", execExists, -2)
	RegisterCommand("flushdb", execFlushDB, -1)
	RegisterCommand("type", execType, 2)
	RegisterCommand("rename", execRename, 3)
	RegisterCommand("renamenx", execRenameNX, 3)
	RegisterCommand("expire", execExpire, 3)
	RegisterCommand("pexpireat", execPexpireat, 3)
}

// DEL
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleteNum := db.Removes(keys...)
	if deleteNum > 0 {
		db.addAof(utils.ToCmdLine3("del", args...))
	}
	return reply.MakeIntReply(int64(deleteNum))
}

// EXISTS
func execExists(db *DB, args [][]byte) resp.Reply {
	var result int64
	for _, arg := range args {
		key := string(arg)
		if _, ok := db.GetEntity(key); ok {
			result++
		}
	}
	return reply.MakeIntReply(result)
}

// FLUSH
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	db.addAof(utils.ToCmdLine3("flushdb", args...))
	return reply.MakeOkReply()
}

// TYPE
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	if entity, ok := db.GetEntity(key); !ok {
		return reply.MakeStatusReply("none")
	} else {
		switch entity.Data.(type) {
		case []byte:
			return reply.MakeStatusReply("string")
		}
		//TODO:
	}
	return &reply.UnknowErrorReply{}

}

// RENAME
func execRename(db *DB, args [][]byte) resp.Reply {
	oldName := string(args[0])
	newName := string(args[1])
	if value, ok := db.GetEntity(oldName); !ok {
		return reply.MakeErrReply("no sush key")
	} else {
		db.PutEntity(newName, value)
		db.Remove(oldName)
		db.addAof(utils.ToCmdLine3("rename", args...))
	}
	return reply.MakeOkReply()
}

// RENAMENX
func execRenameNX(db *DB, args [][]byte) resp.Reply {
	oldName := string(args[0])
	newName := string(args[1])
	if _, ok := db.GetEntity(newName); ok {
		return reply.MakeIntReply(0)
	}
	if value, ok := db.GetEntity(oldName); !ok {
		return reply.MakeErrReply("no sush key")
	} else {
		db.PutEntity(newName, value)
		db.Remove(oldName)
		db.addAof(utils.ToCmdLine3("renamenx", args...))
	}
	return reply.MakeIntReply(1)
}

func execExpire(db *DB, args [][]byte) resp.Reply {
	var (
		key string
		ttl int64
		err error
	)
	key = string(args[0])
	if ttl, err = strconv.ParseInt(string(args[1]), 10, 64); err != nil {
		return reply.MakeErrReply("invalid expire time in set")
	}
	d := time.Duration(ttl) * time.Second
	if _, ok := db.GetEntity(key); !ok {
		return reply.MakeIntReply(0)
	}
	expireAt := time.Now().Add(d)
	db.Expire(key, expireAt)
	return reply.MakeIntReply(1)
}

func execPexpireat(db *DB, args [][]byte) resp.Reply {
	var (
		key string
		raw int64
		err error
	)
	key = string(args[0])
	raw, err = strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return reply.MakeErrReply("ERR value is not an integer or out of range")
	}

	expireAt := time.Unix(0, raw*int64(time.Millisecond))
	if _, ok := db.GetEntity(key); !ok {
		return reply.MakeIntReply(0)
	}
	db.Expire(key, expireAt)
	db.addAof(aof.MakeExpireCmd(key, expireAt))
	return reply.MakeIntReply(1)
}

// KEYS
// func execKeys(db *DB, args [][]byte) resp.Reply {}
