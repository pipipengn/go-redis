package database

import (
	"go-redis/resp/interface"
	"go-redis/resp/reply"
)

type TestDB struct {
}

func NewTestDB() *TestDB {
	return &TestDB{}
}

func (t *TestDB) Exec(client respinterface.Connection, args CmdLine) respinterface.Reply {
	return reply.NewMultiBulkReply(args)
}

func (t *TestDB) Close() {

}

func (t *TestDB) AfterClientClose(client respinterface.Connection) {

}
