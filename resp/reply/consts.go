package reply

type PongReply struct {
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

func (reply *PongReply) ToByte() []byte {
	return []byte("+PONG\r\n")
}

type OkReply struct {
}

func MakeOkReply() *OkReply {
	return &OkReply{}
}

func (reply *OkReply) ToByte() []byte {
	return []byte("+ok\r\n")
}

// 空字符串
type NullBulkReply struct {
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

func (reply *NullBulkReply) ToByte() []byte {
	return []byte("$-1\r\n")
}

// 空数组
type EmptyMultiBulkReply struct {
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

func (reply *EmptyMultiBulkReply) ToByte() []byte {
	return []byte("*0\r\n")
}

type NoReply struct {
}

func MakeNoReply() *NoReply {
	return &NoReply{}
}

func (reply *NoReply) ToByte() []byte {
	return []byte("")
}
