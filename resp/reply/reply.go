package reply

import (
	"bytes"
	"go_redis/interface/resp"
	"strconv"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

// 将Args转化成符合redis协议的字符串回复
type BulkReply struct {
	Args []byte
}

func MakeBulkReply(args []byte) *BulkReply {
	return &BulkReply{Args: args}
}

func (reply *BulkReply) ToByte() []byte {
	if len(reply.Args) == 0 {
		return nullBulkReplyBytes
	}
	return []byte("$" + strconv.Itoa(len(reply.Args)) + CRLF + string(reply.Args) + CRLF)
}

type ErrorReply interface {
	Error() string
	ToByte() []byte
}

// 将多个字符串转化成符合redis协议的回复
type MultiBulkReply struct {
	Args [][]byte
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{Args: args}
}

func (reply *MultiBulkReply) ToByte() []byte {
	argLen := len(reply.Args)
	var buf bytes.Buffer
	buf.WriteString("*" + strconv.Itoa(argLen) + CRLF)
	for _, arg := range reply.Args {
		if arg == nil {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

type StatusReply struct {
	Status string
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{Status: status}
}

func (reply *StatusReply) ToByte() []byte {
	return []byte("+" + reply.Status + CRLF)
}

type IntReply struct {
	Code int64
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

func (reply *IntReply) ToByte() []byte {
	return []byte(":" + strconv.FormatInt(reply.Code, 10) + CRLF)
}

type StandardErrReply struct {
	Status string
}

func (reply *StandardErrReply) ToByte() []byte {
	return []byte("-" + reply.Status + CRLF)
}

func (reply *StandardErrReply) Error() string {
	return reply.Status
}

func MakeErrReply(status string) *StandardErrReply {
	return &StandardErrReply{
		Status: status,
	}
}

func IsErrorReply(reply resp.Reply) bool {
	return reply.ToByte()[0] == '-'
}
