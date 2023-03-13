package parser

import (
	"bufio"
	"errors"
	"go_redis/interface/resp"
	"go_redis/lib/logger"
	"go_redis/resp/reply"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
)

type PayLoad struct {
	Data resp.Reply
	Err  error
}

type ReadState struct {
	readMultiLine     bool //是否解析多行数据
	expectedArgsCount int  //目标解析参数数量
	msgType           byte
	args              [][]byte //解析出来的参数
	bulkLen           int64    //数据块的总长度
}

func (state *ReadState) finished() bool {
	return state.expectedArgsCount > 0 && len(state.args) == state.expectedArgsCount
}

func ParserStream(reader io.Reader) <-chan *PayLoad {
	ch := make(chan *PayLoad)
	go parser0(reader, ch)
	return ch
}

func parser0(reader io.Reader, ch chan *PayLoad) {
	var (
		state ReadState
		err   error
		msg   []byte
	)
	defer func() {
		if err := recover(); err != nil {
			logger.Error("ErrorInfo: ", err, "\r", string(debug.Stack()))
		}
	}()
	bufReader := bufio.NewReader(reader)
	for {
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			if ioErr {
				ch <- &PayLoad{
					Err: err,
				}
				close(ch)
				return
			}
			ch <- &PayLoad{
				Err: err,
			}
			state = ReadState{}
			continue
		}
		//判断是否为多行解析模式
		if !state.readMultiLine {
			//*3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nVALUE\r\n
			if msg[0] == '*' {
				//数组模式
				if err = parseMultiBulkHeader(msg, &state); err != nil {
					ch <- &PayLoad{
						Err: err,
					}
					state = ReadState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					//发送了空命令
					ch <- &PayLoad{
						Data: &reply.EmptyMultiBulkReply{},
					}
					state = ReadState{}
					continue
				}
			} else if msg[0] == '$' {
				//多行字符串
				if err = parseBulkHeader(msg, &state); err != nil {
					ch <- &PayLoad{
						Err: err,
					}
					state = ReadState{}
					continue
				}
				if state.bulkLen == -1 {
					ch <- &PayLoad{
						Data: &reply.NullBulkReply{},
					}
					state = ReadState{}
					continue
				}
			} else {
				// + - :
				var result resp.Reply
				result, err = parseSingleLineReply(msg)
				ch <- &PayLoad{
					Data: result,
					Err:  err,
				}
				state = ReadState{}
				continue
			}
		} else {
			if err = readBody(msg, &state); err != nil {
				ch <- &PayLoad{
					Err: err,
				}
			}
			if state.finished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.MakeMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.MakeBulkReply(state.args[0])
				}
				ch <- &PayLoad{
					Data: result,
				}
				state = ReadState{}
			}
		}
	}
}

// 读取一行数据
func readLine(bufReader *bufio.Reader, state *ReadState) (msg []byte, isIOErr bool, err error) {
	//两种情况
	//1.普通类型 按照\r\n切分
	if state.bulkLen == 0 {
		//\r\n切分
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
	} else {
		//2.Bulk Strings 严格按照$后的数字进行切分
		msg = make([]byte, int(state.bulkLen)+2)
		if _, err = io.ReadFull(bufReader, msg); err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

// *3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nVALUE\r\n
func readBody(msg []byte, state *ReadState) error {
	line := msg[0 : len(msg)-2]
	var err error
	if line[0] == '$' {
		// bulk reply
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error: " + string(msg))
		}
		if state.bulkLen <= 0 { // null bulk in multi bulks
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}

// *数组解析
func parseMultiBulkHeader(msg []byte, state *ReadState) (err error) {
	var expectedLine uint64
	expectedLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil {
		return errors.New("protocol error:" + string(msg))
	}
	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return
	} else if expectedLine > 0 {
		state.msgType = msg[0]
		state.expectedArgsCount = int(expectedLine)
		state.readMultiLine = true
		state.args = make([][]byte, 0, expectedLine)
		return
	} else {
		return errors.New("protocol error:" + string(msg))
	}
}

// $多行字符串解析
func parseBulkHeader(msg []byte, state *ReadState) (err error) {
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error:" + string(msg))
	}
	if state.bulkLen == -1 {
		return nil
	} else if state.bulkLen > 0 {
		state.expectedArgsCount = 1
		state.msgType = msg[0]
		state.readMultiLine = true
		state.args = make([][]byte, 0, 1)
		return
	} else {
		return errors.New("protocol error:" + string(msg))
	}
}

func parseSingleLineReply(msg []byte) (resp resp.Reply, err error) {
	strMsg := strings.TrimSuffix(string(msg), "\r\n")
	switch msg[0] {
	case '+':
		resp = reply.MakeStatusReply(strMsg[1:])
	case '-':
		resp = reply.MakeErrReply(strMsg[1:])
	case ':':
		var val int64
		val, err = strconv.ParseInt(strMsg[1:], 10, 64)
		if err != nil {
			return nil, errors.New("protocol error:" + string(msg))
		}
		resp = reply.MakeIntReply(val)
	}
	return
}
