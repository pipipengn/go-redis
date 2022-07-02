package database

import (
	"go-redis/core/dict"
	respinterface "go-redis/resp/interface"
	"go-redis/resp/reply"
	"strings"
)

type DB struct {
	index      int
	dict       dict.Interface
	addAofFunc func([][]byte)
}

func NewDB(dict dict.Interface) *DB {
	return &DB{
		dict:       dict,
		addAofFunc: func([][]byte) {},
	}
}

func (db *DB) Exec(client respinterface.Connection, cmd [][]byte) respinterface.Reply {
	action := strings.ToLower(string(cmd[0]))
	command, ok := cmdTable[action]
	if !ok {
		return reply.NewErrReply("ERR unknow command: " + action)
	}
	if !validateArgNUm(command.argNum, cmd) {
		return reply.NewArgNumErrReply(action)
	}
	return command.executor(db, cmd[1:])
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

func (db *DB) GetEntity(key string) (*dict.DataEntity, bool) {
	raw, ok := db.dict.Get(key)
	if !ok {
		return nil, false
	}
	return &dict.DataEntity{Data: raw}, true
}

func (db *DB) SetEntity(key string, entity *dict.DataEntity) (rowAffected int) {
	return db.dict.Set(key, entity.Data)
}

func (db *DB) SetIfExists(key string, entity *dict.DataEntity) (rowAffected int) {
	return db.dict.SetIfExists(key, entity.Data)
}

func (db *DB) SetIfAbsent(key string, entity *dict.DataEntity) (rowAffected int) {
	return db.dict.SetIfAbsent(key, entity.Data)
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
