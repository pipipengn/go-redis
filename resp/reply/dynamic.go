package reply

import (
	"bytes"
	"go-redis/resp/interface"
	"strconv"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

// BulkReply ========================================================================
type BulkReply struct {
	Arg []byte
}

func NewBulkReply(arg []byte) *BulkReply {
	return &BulkReply{Arg: arg}
}

func (r *BulkReply) ToBytes() []byte {
	if len(r.Arg) == 0 {
		return []byte(string(nullBulkReplyBytes) + CRLF)
	}
	return []byte("$" + strconv.Itoa(len(r.Arg)) + CRLF + string(r.Arg) + CRLF)
}

// MultiBulkReply ========================================================================
type MultiBulkReply struct {
	Args [][]byte
}

func NewMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{Args: args}
}

func (r *MultiBulkReply) ToBytes() []byte {
	buf := bytes.Buffer{}
	buf.WriteString("*" + strconv.Itoa(len(r.Args)) + CRLF)
	for _, arg := range r.Args {
		if arg == nil {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

// StatusReply ========================================================================
type StatusReply struct {
	Status string
}

func NewStatusReply(status string) *StatusReply {
	return &StatusReply{Status: status}
}

func (r *StatusReply) ToBytes() []byte {
	return []byte("+" + r.Status + CRLF)
}

// IntReply ========================================================================
type IntReply struct {
	Code int
}

func NewIntReply(code int) *IntReply {
	return &IntReply{Code: code}
}

func (r *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(int64(r.Code), 10) + CRLF)
}

// ErrReply ===============================================================
type ErrReply struct {
	Status string
}

func NewErrReply(status string) *ErrReply {
	return &ErrReply{Status: status}
}

func (r *ErrReply) Error() string {
	return r.Status
}

func (r *ErrReply) ToBytes() []byte {
	return []byte("-" + r.Status + CRLF)
}

func IsErrorReply(reply respinterface.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
