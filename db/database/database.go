package database

import (
	"go-redis/db/dict"
	respinterface "go-redis/resp/interface"
	"go-redis/resp/reply"
	"strings"
)

type DB struct {
	index int
	dict  dict.Interface
}

func NewDB(dict dict.Interface) *DB {
	return &DB{dict: dict}
}

type ExecFunc func(db *DB, args [][]byte) respinterface.Reply

func (db *DB) Exec(client respinterface.Connection, cmd [][]byte) respinterface.Reply {
	action := strings.ToLower(string(cmd[0]))
	command, ok := cmdTable[action]
	if !ok {
		return reply.NewErrReply("ERR unknow command: " + action)
	}
	if !validateArgNUm(command.argNum, cmd) {
		return reply.NewArgNumErrReply(action)
	}
	// TODO ExecFunc Impl
	return command.executor(db, cmd[1:])
}

func (db *DB) Close() {
	//TODO implement me
	panic("implement me")
}

func (db *DB) AfterClientClose(client respinterface.Connection) {
	//TODO implement me
	panic("implement me")
}

// EXISTS k1 k2 k3 = -2
// SET key val = 3
func validateArgNUm(argNum int, cmd [][]byte) bool {
	cmdLen := len(cmd)
	if argNum >= 0 {
		return argNum == cmdLen
	}
	return cmdLen >= -argNum
}

// ===================

func (db *DB) Get(key string) (*dict.DataEntity, bool) {
	raw, ok := db.dict.Get(key)
	if !ok {
		return nil, false
	}
	return &dict.DataEntity{Data: raw}, true
}

func (db *DB) SetEntity(key string, entity *dict.DataEntity) int {
	return db.dict.Set(key, entity)
}

func (db *DB) SetIfExists(key string, entity *dict.DataEntity) int {
	return db.dict.SetIfExists(key, entity)
}

func (db *DB) SetIfAbsent(key string, entity *dict.DataEntity) int {
	return db.dict.SetIfAbsent(key, entity)
}

func (db *DB) Remove(key string) {
	db.dict.Remove(key)
}

func (db *DB) Removes(keys ...string) (rowAffected int) {
	rowAffected = 0
	for _, key := range keys {
		if affected := db.dict.Remove(key); affected == 1 {
			rowAffected++
		}
	}
	return rowAffected
}

func (db *DB) Flush() {
	db.dict.Clear()
}
