package aof

import (
	"go_redis/config"
	databaseface "go_redis/interface/database"
	"go_redis/lib/logger"
	"go_redis/lib/utils"
	"go_redis/resp/reply"
	"os"
	"strconv"
)

type Cmdline = [][]byte

type payload struct {
	cmdLine Cmdline
	dbIndex int
}

const aofBufferSize = 1 << 16

type AofHandler struct {
	database    databaseface.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFileName string
	currentDB   int
}

func NewAofHandler(database databaseface.Database) (handler *AofHandler, err error) {
	fileName := config.Properties.AppendFilename
	handler = &AofHandler{
		aofFileName: fileName,
		database:    database,
		aofChan:     make(chan *payload, aofBufferSize),
	}
	handler.LoadAof()
	go func() {
		handler.handleAof()
	}()
	aofile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofile

	return
}

func (aof *AofHandler) AddAof(dbIndex int, cmd Cmdline) {
	if config.Properties.AppendOnly && aof.aofChan != nil {
		aof.aofChan <- &payload{
			cmdLine: cmd,
			dbIndex: dbIndex,
		}
	}
}

// <-aofChan 落盘
func (aof *AofHandler) handleAof() {
	aof.currentDB = 0
	var (
		data []byte
		err  error
	)
	for p := range aof.aofChan {
		if p.dbIndex != aof.currentDB {
			data = reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))).ToByte()
			if _, err = aof.aofFile.Write(data); err != nil {
				logger.Error(err)
				continue
			}
			aof.currentDB = p.dbIndex
		}
		data = reply.MakeMultiBulkReply(p.cmdLine).ToByte()
		if _, err = aof.aofFile.Write(data); err != nil {
			logger.Error(err)
			continue
		}
	}
}

func (aof *AofHandler) LoadAof() {

}
