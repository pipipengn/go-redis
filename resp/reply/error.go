package reply

// UnknowErrReply unknow error ===============================================================
type UnknowErrReply struct{}

var unknowErrReply = new(UnknowErrReply)

func NewUnknowErrReply() *UnknowErrReply {
	return unknowErrReply
}

func (r *UnknowErrReply) Error() string {
	return "Err unknow"
}

func (r *UnknowErrReply) ToBytes() []byte {
	return []byte("-Err unknow\r\n")
}

// ArgNumErrReply arg number error ===============================================================
type ArgNumErrReply struct {
	cmd string
}

func NewArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{cmd: cmd}
}

func (r *ArgNumErrReply) Error() string {
	return "ERR wrong number of arguments for '" + r.cmd + "' command"
}

func (r *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + r.cmd + "' command\r\n")
}

// SyntaxErrReply syntax error ===============================================================
type SyntaxErrReply struct{}

var syntaxErrReply = new(SyntaxErrReply)

func NewSyntaxErrReply() *SyntaxErrReply {
	return syntaxErrReply
}

func (r *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (r *SyntaxErrReply) ToBytes() []byte {
	return []byte("-Err syntax error\r\n")
}

// WrongTypeErrReply wrong type error ===============================================================
type WrongTypeErrReply struct{}

var wrongTypeErrReply = new(WrongTypeErrReply)

func NewWrongTypeErrReply() *WrongTypeErrReply {
	return wrongTypeErrReply
}

func (r *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

func (r *WrongTypeErrReply) ToBytes() []byte {
	return []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")
}

// ProtocolErrReply protocol error ===============================================================
type ProtocolErrReply struct {
	msg string
}

func NewProtocolErrReply(msg string) *ProtocolErrReply {
	return &ProtocolErrReply{msg: msg}
}

func (r *ProtocolErrReply) Error() string {
	return "ERR Protocol error: '" + r.msg
}

func (r *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + r.msg + "'\r\n")
}
