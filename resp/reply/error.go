package reply

type UnknowErrorReply struct {
}

var unknowErrBytes = []byte("-Err unknown\r\n")

func (errReply *UnknowErrorReply) Error() string {
	return "Err unknown"
}

func (errReply *UnknowErrorReply) ToByte() []byte {
	return unknowErrBytes
}

// 参数数量错误
type ArgNumErrReply struct {
	Cmd string
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

func (errReply *ArgNumErrReply) Error() string {
	return "ERR wrong number of arguments for '" + errReply.Cmd + "' command"
}

func (errReply *ArgNumErrReply) ToByte() []byte {
	return []byte("-ERR wrong number of arguments for '" + errReply.Cmd + "' command\r\n")
}

// 语法错误
type SyntaxErrReply struct{}

var syntaxErrBytes = []byte("-Err syntax error\r\n")
var theSyntaxErrReply = &SyntaxErrReply{}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

func (r *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

func (r *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

// 数据类型错误
type WrongTypeErrReply struct{}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")

func (r *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

func (r *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

//接口协议错误 不符合redis协议

type ProtocolErrReply struct {
	Msg string
}

func (r *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + r.Msg + "'\r\n")
}

func (r *ProtocolErrReply) Error() string {
	return "ERR Protocol error: '" + r.Msg
}
