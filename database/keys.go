package database

import (
	"go_redis/interface/resp"
	"go_redis/resp/reply"
)

func init() {
	RegisterCommand("del", execDel, -2)
	RegisterCommand("exists", execExists, -2)
	RegisterCommand("flushdb", execFlushDB, -1)
	RegisterCommand("type", execType, 2)
	RegisterCommand("rename", execRename, 3)
	RegisterCommand("renamenx", execRenameNX, 3)
}

// DEL
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleteNum := db.Removes(keys...)
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
	}
	return reply.MakeIntReply(1)
}

// KEYS
// func execKeys(db *DB, args [][]byte) resp.Reply {}
