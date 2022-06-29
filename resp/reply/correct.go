package reply

// PongReply pong ========================================================
type PongReply struct {
}

var pongReply = new(PongReply)

func NewPongReply() *PongReply {
	return pongReply
}

func (r *PongReply) ToBytes() []byte {
	return []byte("+PONG\r\n")
}

// OkReply ok ========================================================
type OkReply struct{}

var okReply = new(OkReply)

func NewOkReply() *OkReply {
	return okReply
}

func (r *OkReply) ToBytes() []byte {
	return []byte("+OK\r\n")
}

// EmptyBulkReply empty string ========================================================
type EmptyBulkReply struct{}

var emptyBulkReply = new(EmptyBulkReply)

func NewEmptyBulkReply() *EmptyBulkReply {
	return emptyBulkReply
}

func (r *EmptyBulkReply) ToBytes() []byte {
	return []byte("$-1\r\n")
}

// EmptyMultiBulkReply empty array ========================================================
type EmptyMultiBulkReply struct{}

var emptyMultiBulkReply = new(EmptyMultiBulkReply)

func NewEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return emptyMultiBulkReply
}

func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return []byte("*0\r\n")
}

// NoReply return nil ========================================================
type NoReply struct{}

var noReply = new(NoReply)

func NewNoReply() *NoReply {
	return noReply
}

func (r *NoReply) ToBytes() []byte {
	return []byte("")
}
