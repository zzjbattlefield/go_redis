package aof

import (
	"context"
	databaseface "go_redis/interface/database"
	"go_redis/lib/logger"
	"go_redis/lib/utils"
	"go_redis/resp/connection"
	"go_redis/resp/parser"
	"go_redis/resp/reply"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	FsyncAlways      = "always"
	FsyncEverySecond = "everysec"
	FsyncNo          = "no"
)

type payLoad struct {
	dbIndex int
	cmdLine [][]byte
}

type Persister struct {
	fileName       string
	db             databaseface.Database
	aofChan        chan *payLoad
	aofFile        *os.File
	currentDB      int
	aofFinshedChan chan struct{}
	ctx            context.Context
	cancel         context.CancelFunc
	aofFsync       string
	pausingAof     sync.Mutex
}

func NewPersister(db databaseface.Database, fileName string, load bool, fsyncType string) (*Persister, error) {
	p := &Persister{}
	p.fileName = fileName
	p.db = db
	p.aofFsync = strings.ToLower(fsyncType)
	if load {
		p.loadAof()
	}
	aofFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	go func() {
		p.listenCmd()
	}()
	p.aofFile = aofFile
	p.aofChan = make(chan *payLoad, 1<<16)
	p.aofFinshedChan = make(chan struct{})
	p.ctx, p.cancel = context.WithCancel(context.Background())
	if p.aofFsync == FsyncEverySecond {
		go p.fsyncEverySecond()
	}
	return p, nil
}

// 从aof文件恢复数据 要在listenCmd之前调用
func (p *Persister) loadAof() {
	aofChan := p.aofChan
	//load的时候不能有aofchan 不然load的数据又回写入aof里
	if aofChan != nil {
		p.aofChan = nil
		defer func(aofChan chan *payLoad) {
			p.aofChan = aofChan
		}(aofChan)
	}
	f, err := os.Open(p.fileName)
	if err != nil {
		logger.Warn("loadAof open file err:" + err.Error())
	}
	defer f.Close()
	ch := parser.ParserStream(f)
	fakeConn := &connection.Connection{}
	for payLoad := range ch {
		if payLoad.Err != nil {
			if payLoad.Err == io.EOF {
				break
			}
			logger.Error("ParserStream error:%s", err.Error())
			continue
		}
		if payLoad.Data == nil {
			logger.Error("empty payload")
			break
		}
		r, ok := payLoad.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk protocol")
			continue
		}
		ret := p.db.Exec(fakeConn, r.Args)
		if reply.IsErrorReply(ret) {
			logger.Error("loadAof error:%s", string(ret.ToByte()))
		}
		if strings.ToLower(string(r.Args[0])) == "select" {
			dbIndex, _ := strconv.Atoi(string(r.Args[1]))
			p.currentDB = dbIndex
		}
	}
}

func (p *Persister) listenCmd() {
	logger.Info("aof listenCmd start")
	for payload := range p.aofChan {
		p.writeAof(payload)
	}
	p.aofFinshedChan <- struct{}{}
}

func (p *Persister) SaveCmdLine(index int, cmdLine [][]byte) {
	if p.aofChan == nil {
		return
	}
	payLoad := &payLoad{
		dbIndex: index,
		cmdLine: cmdLine,
	}
	if p.aofFsync == FsyncAlways {
		p.writeAof(payLoad)
	}
	p.aofChan <- payLoad
}

func (p *Persister) writeAof(payload *payLoad) {
	p.pausingAof.Lock()
	defer p.pausingAof.Unlock()
	if payload.dbIndex != p.currentDB {
		//切换数据库
		selectCmd := utils.ToCmdLine("SELECT" + strconv.Itoa(payload.dbIndex))
		data := reply.MakeMultiBulkReply(selectCmd).ToByte()
		if _, err := p.aofFile.Write(data); err != nil {
			logger.Warn(err)
			return //skip
		}
	}
	data := reply.MakeMultiBulkReply(payload.cmdLine).ToByte()
	if _, err := p.aofFile.Write(data); err != nil {
		logger.Warn(err)
		return //skip
	}
	if p.aofFsync == FsyncAlways {
		p.aofFile.Sync()
	}
}

func (p *Persister) fsyncEverySecond() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			p.pausingAof.Lock()
			if err := p.aofFile.Sync(); err != nil {
				logger.Error("fsync aof file error: %v", err)
			}
			p.pausingAof.Unlock()
		case <-p.ctx.Done():
			return
		}
	}
}
