package database

import (
	"go_redis/interface/resp"
	"go_redis/lib/logger"
	"go_redis/resp/reply"
	"strconv"
	"strings"
)

type Database struct {
	dbSet []*DB
}

func NewDatabase() *Database {
	database := &Database{
		dbSet: make([]*DB, 16),
	}
	for index := range database.dbSet {
		db := MakeDB()
		db.index = index
		database.dbSet[index] = db
	}
	return database
}

func (database *Database) Exec(client resp.Connect, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			reply.MakeArgNumErrReply("select")
		}
		return execSelect(client, database, args[1:])
	}
	dbIndex := client.GetDBIndex()
	db := database.dbSet[dbIndex]
	return db.Exec(client, args)
}
func (database *Database) Close() {

}
func (database *Database) AfterClientClose(c resp.Connect) {

}

func execSelect(c resp.Connect, database *Database, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}
	if dbIndex > len(database.dbSet) {
		return reply.MakeErrReply("ERR invalid DB index")
	}
	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}
