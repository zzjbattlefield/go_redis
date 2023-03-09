package resp

type Connect interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
