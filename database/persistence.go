package database

import "go_redis/aof"

func (database *Database) bindPersister(aofHandler *aof.Persister) {
	database.persister = aofHandler
	for _, db := range database.dbSet {
		index := db.index
		db.addAof = func(line CmdLine) {
			aofHandler.SaveCmdLine(index, line)
		}
	}
}
