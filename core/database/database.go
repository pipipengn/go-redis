package database

import (
	"go-redis/config"
	"go-redis/core/aof"
	"go-redis/core/dict"
	respinterface "go-redis/resp/interface"
	"go-redis/resp/reply"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type Database struct {
	dbSet      []*DB
	aofHandler *aof.Handler
}

func New() *Database {
	if config.Config.Databases <= 0 {
		config.Config.Databases = 16
	}
	dbs := make([]*DB, config.Config.Databases)
	for i := range dbs {
		db := NewDB(dict.NewSyncDict())
		db.index = i
		dbs[i] = db
	}
	database := &Database{dbSet: dbs}
	if config.Config.AppendOnly {
		aofHandler, err := aof.NewHandler(database)
		if err != nil {
			zap.S().Fatalf("cannot init aof: %v", zap.Error(err))
		}
		database.aofHandler = aofHandler

		for _, db := range database.dbSet {
			db2 := db
			db2.addAofFunc = func(cmd [][]byte) {
				database.aofHandler.AddAof(db2.index, cmd)
			}
		}
	}

	return database
}

// Exec set k v / get k / select 1
func (d *Database) Exec(client respinterface.Connection, args [][]byte) respinterface.Reply {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Error(err)
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.NewArgNumErrReply(cmdName)
		}
		return d.ExecSelect(client, args[1:])
	}
	return d.dbSet[client.GetDBIndex()].Exec(client, args)
}

func (d *Database) Close() {

}

func (d *Database) AfterClientClose(client respinterface.Connection) {

}

// ExecSelect select 1
func (d *Database) ExecSelect(c respinterface.Connection, args [][]byte) respinterface.Reply {
	dbIdx, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.NewErrReply("ERR invalid DB index")
	}
	if dbIdx >= len(d.dbSet) {
		return reply.NewErrReply("ERR DB index is out of range")
	}
	c.SelectDB(dbIdx)
	return reply.NewOkReply()
}
