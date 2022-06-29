package reply

import (
	"bytes"
	"go-redis/resp"
	"strconv"
)

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

// BulkReply ========================================================================
type BulkReply struct {
	arg []byte
}

func NewBulkReply(arg []byte) *BulkReply {
	return &BulkReply{arg: arg}
}

func (r *BulkReply) ToBytes() []byte {
	if len(r.arg) == 0 {
		return []byte(string(nullBulkReplyBytes) + CRLF)
	}
	return []byte("$" + strconv.Itoa(len(r.arg)) + CRLF + string(r.arg) + CRLF)
}

// MultiBulkReply ========================================================================
type MultiBulkReply struct {
	args [][]byte
}

func NewMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{args: args}
}

func (r *MultiBulkReply) ToBytes() []byte {
	buf := bytes.Buffer{}
	buf.WriteString("*" + strconv.Itoa(len(r.args)) + CRLF)
	for _, arg := range r.args {
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
	status string
}

func NewStatusReply(status string) *StatusReply {
	return &StatusReply{status: status}
}

func (r *StatusReply) ToBytes() []byte {
	return []byte("+" + r.status + CRLF)
}

// IntReply ========================================================================
type IntReply struct {
	code int
}

func NewIntReply(code int) *IntReply {
	return &IntReply{code: code}
}

func (r *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(int64(r.code), 10) + CRLF)
}

// CustomeErrReply ===============================================================
type CustomeErrReply struct {
	status string
}

func NewCustomeErrReply(status string) *CustomeErrReply {
	return &CustomeErrReply{status: status}
}

func (r *CustomeErrReply) Error() string {
	return r.status
}

func (r *CustomeErrReply) ToBytes() []byte {
	return []byte("-" + r.status + CRLF)
}

func IsErrorReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
