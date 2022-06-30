package database

import (
	"go-redis/resp/interface"
)

type CmdLine [][]byte

type Interface interface {
	Exec(client respinterface.Connection, args CmdLine) respinterface.Reply
	Close()
	AfterClientClose(client respinterface.Connection)
}

type DataEntity struct {
	Data interface{}
}
