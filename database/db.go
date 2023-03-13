package database

import (
	"go_redis/datastruct/dict"
	"go_redis/interface/database"
	"go_redis/interface/resp"
	"go_redis/resp/reply"
	"strings"
)

type DB struct {
	index int
	data  dict.Dict
}

func MakeDB() *DB {
	return &DB{
		data: dict.MakeSyncDict(),
	}
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func (db *DB) Exec(c resp.Connect, cmdLine CmdLine) resp.Reply {
	// ping set delete
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknow command")
	}
	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	fun := cmd.exector
	return fun(db, cmdLine[1:])
}

// 参数校验 如何是不定长参数=-2
func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		//定长参数
		return argNum == arity
	}
	return argNum >= -arity
}

func (db *DB) GetEntity(key string) (data *database.DataEntity, ok bool) {
	var raw any
	raw, ok = db.data.Get(key)
	if !ok {
		return
	}
	data, ok = raw.(*database.DataEntity)
	return
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

func (db *DB) Remove(key string) int {
	return db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) (deleted int) {
	for _, item := range keys {
		if _, exist := db.data.Get(item); exist {
			deleted++
			db.Remove(item)
		}
	}
	return
}

func (db *DB) Flush() {
	db.data.Clear()
}
