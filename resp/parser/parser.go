package parser

import (
	"bufio"
	"errors"
	"go-redis/resp/interface"
	"go-redis/resp/reply"
	"go.uber.org/zap"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
)

type DataStream struct {
	Data respinterface.Reply
	Err  error
}

type readState struct {
	readingMultiLine bool
	expectArgsCount  int
	msgType          byte
	args             [][]byte
	bulkLen          int
}

func (r *readState) finished() bool {
	return r.expectArgsCount > 0 && r.expectArgsCount == len(r.args)
}

// ParseStream tcp layer call
func ParseStream(conn io.Reader) <-chan *DataStream {
	ch := make(chan *DataStream)
	go parse(conn, ch)
	return ch
}

// concurrent call
func parse(conn io.Reader, ch chan *DataStream) {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Error(string(debug.Stack()))
		}
	}()

	bufReader := bufio.NewReader(conn)
	var state readState
	for {
		msg, isIOErr, err := readLine(bufReader, &state)
		if err != nil {
			if isIOErr {
				ch <- &DataStream{Err: err}
				close(ch)
				return
			}
			ch <- &DataStream{Err: err}
			state = readState{}
			continue
		}

		// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
		if !state.readingMultiLine {
			switch msg[0] {
			case '*':
				if err := parseMultiBulkHeader(msg, &state); err != nil {
					ch <- &DataStream{Err: err}
					state = readState{}
					continue
				}
				if state.expectArgsCount == 0 {
					ch <- &DataStream{Data: reply.NewEmptyMultiBulkReply()}
					state = readState{}
					continue
				}
			case '$':
				if err := parseBulkHeader(msg, &state); err != nil {
					ch <- &DataStream{Err: err}
					state = readState{}
					continue
				}
				if state.bulkLen == -1 {
					ch <- &DataStream{Data: reply.NewEmptyBulkReply()}
					state = readState{}
					continue
				}
			case '+', '-', ':':
				singleLineReply, err := parseSingleLineReply(msg)
				ch <- &DataStream{Data: singleLineReply, Err: err}
				state = readState{}
				continue
			}
		} else {
			if err := readBody(msg, &state); err != nil {
				ch <- &DataStream{Err: err}
				state = readState{}
				continue
			}
			if state.finished() {
				var result respinterface.Reply
				if state.msgType == '*' {
					result = reply.NewMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.NewBulkReply(state.args[0])
				}

				ch <- &DataStream{Data: result}
				state = readState{}
			}
		}
	}
}

// *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
func readLine(reader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error
	if state.bulkLen == 0 {
		msg, err = reader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error1: " + string(msg))
		}
	} else {
		msg = make([]byte, state.bulkLen+2)
		if _, err := io.ReadFull(reader, msg); err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error2: " + string(msg))
		}
		state.bulkLen = 0
	}

	return msg, false, nil
}

// *3\r\n
func parseMultiBulkHeader(msg []byte, state *readState) error {
	expectedLine, err := strconv.Atoi(string(msg[1 : len(msg)-2]))
	if err != nil {
		return errors.New("protocol error3: " + string(msg))
	}

	if expectedLine == 0 {
		state.expectArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.args = make([][]byte, 0, 3)
		state.expectArgsCount = expectedLine
		return nil
	} else {
		return errors.New("protocol error4: " + string(msg))
	}
}

// $3\r\n
func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.Atoi(string(msg[1 : len(msg)-2]))
	if err != nil {
		return errors.New("protocol error5: " + string(msg))
	}

	if state.bulkLen == -1 {
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("protocol error6: " + string(msg))
	}
}

// +OK\r\n  -err\r\n  :4\r\n
func parseSingleLineReply(msg []byte) (respinterface.Reply, error) {
	s := strings.TrimSuffix(string(msg), "\r\n")
	var result respinterface.Reply
	switch msg[0] {
	case '+':
		result = reply.NewStatusReply(s[1:])
	case '-':
		result = reply.NewErrReply(s[1:])
	case ':':
		val, err := strconv.Atoi(s[1:])
		if err != nil {
			return nil, errors.New("protocol error7: " + string(msg))
		}
		result = reply.NewIntReply(val)
	}
	return result, nil
}

// $3\r\n
// PING\r\n
func readBody(msg []byte, state *readState) error {
	line := msg[:len(msg)-2]

	if line[0] == '$' {
		var err error
		if state.bulkLen, err = strconv.Atoi(string(line[1:])); err != nil {
			return errors.New("protocol error8: " + string(msg))
		}
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}

	return nil
}
