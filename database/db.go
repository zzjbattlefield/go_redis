package database

import (
	"go_redis/datastruct/dict"
	"go_redis/interface/database"
	"go_redis/interface/resp"
	"go_redis/lib/logger"
	"go_redis/lib/timewheel"
	"go_redis/resp/reply"
	"strings"
	"time"
)

type DB struct {
	index  int
	data   dict.Dict
	ttlMap dict.Dict
	addAof func(CmdLine)
}

func MakeDB() *DB {
	return &DB{
		data:   dict.MakeSyncDict(),
		ttlMap: dict.MakeSyncDict(),
		addAof: func(cmdline CmdLine) {},
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
	raw, exists := db.data.Get(key)
	if !exists || db.isExpired(key) {
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

func (db *DB) Expire(key string, expireTime time.Time) {
	db.ttlMap.Put(key, expireTime)
	timewheel.At(expireTime, key, func() {
		logger.Info("expire key:", key)
		rawExpireTime, ok := db.ttlMap.Get(key)
		if !ok {
			return
		}
		expireTime := rawExpireTime.(time.Time)
		if isExpired := time.Now().After(expireTime); isExpired {
			db.Remove(key)
			logger.Info("expired! remove key:", key)
		}
	})
}

func (db *DB) isExpired(key string) bool {
	rawExpireTime, exist := db.ttlMap.Get(key)
	if !exist {
		return false
	}
	expireTime, _ := rawExpireTime.(time.Time)
	if isExpire := time.Now().After(expireTime); isExpire {
		logger.Info("delete expired key:", key)
		db.data.Remove(key)
	}
	return true
}
