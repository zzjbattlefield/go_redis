package aof

import (
	"strconv"
	"time"
)

func MakeExpireCmd(key string, expiredAt time.Time) (cmdLine [][]byte) {
	cmdLine = make([][]byte, 3)
	cmdLine[0] = []byte("PEXPIREAT")
	cmdLine[1] = []byte(key)
	cmdLine[2] = []byte(strconv.FormatInt(expiredAt.UnixNano()/1e6, 10))
	return
}
