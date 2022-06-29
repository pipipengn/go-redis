package database

import (
	"go-redis/resp"
)

type CmdLine [][]byte

type Interface interface {
	Exec(client resp.Connection, args CmdLine) resp.Reply
	Close()
	AfterClientClose(client resp.Connection)
}

type DataEntity struct {
	Data interface{}
}
