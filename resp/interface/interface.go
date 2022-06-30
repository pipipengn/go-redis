package respinterface

type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}

type Reply interface {
	ToBytes() []byte
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}
