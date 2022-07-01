package database

import (
	"go-redis/core/dict"
	respinterface "go-redis/resp/interface"
	"go-redis/resp/reply"
	"go-redis/utils/wildcard"
)

// Ping ping pong
func Ping(db *DB, args [][]byte) respinterface.Reply {
	return reply.NewPongReply()
}

// ExecDel del k1 k2 k3...
func ExecDel(db *DB, args [][]byte) respinterface.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	rowAffected := db.Removes(keys...)
	return reply.NewIntReply(rowAffected)
}

// ExecExists exists k1 k2 k3
func ExecExists(db *DB, args [][]byte) respinterface.Reply {
	result := 0
	for _, arg := range args {
		key := string(arg)
		if _, ok := db.GetEntity(key); ok {
			result++
		}
	}
	return reply.NewIntReply(result)
}

// ExecFlushDB flush core
func ExecFlushDB(db *DB, args [][]byte) respinterface.Reply {
	db.Flush()
	return reply.NewOkReply()
}

// ExecType type k1
func ExecType(db *DB, args [][]byte) respinterface.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.NewStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.NewStatusReply("string")
	}
	return reply.NewUnknowErrReply()
}

// ExecRename rename k1 k2
func ExecRename(db *DB, args [][]byte) respinterface.Reply {
	key1 := string(args[0])
	key2 := string(args[1])
	entity, exists := db.GetEntity(key1)
	if !exists {
		return reply.NewErrReply("key does not exists")
	}
	db.SetEntity(key2, entity)
	db.Remove(key1)
	return reply.NewOkReply()
}

// ExecRenameNx renamenx k1 k2
func ExecRenameNx(db *DB, args [][]byte) respinterface.Reply {
	key1 := string(args[0])
	key2 := string(args[1])
	entity, exists := db.GetEntity(key1)
	if !exists {
		return reply.NewErrReply("key does not exists")
	}
	rowAffected := db.SetIfAbsent(key2, entity)
	if rowAffected == 1 {
		db.Remove(key1)
	}
	return reply.NewIntReply(rowAffected)
}

// ExecKeys keys *
func ExecKeys(db *DB, args [][]byte) respinterface.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.dict.Range(func(key string, val any) bool {
		if isMatch := pattern.IsMatch(key); isMatch {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.NewMultiBulkReply(result)
}

// ExecGet get key
func ExecGet(db *DB, args [][]byte) respinterface.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.NewEmptyBulkReply()
	}
	return reply.NewBulkReply(entity.Data.([]byte))
}

// ExecSet set k v
func ExecSet(db *DB, args [][]byte) respinterface.Reply {
	key := string(args[0])
	val := args[1]
	db.SetEntity(key, &dict.DataEntity{Data: val})
	return reply.NewOkReply()
}

// ExecSetNx setnx k1 v1
func ExecSetNx(db *DB, args [][]byte) respinterface.Reply {
	key := string(args[0])
	val := args[1]
	rowAffected := db.SetIfAbsent(key, &dict.DataEntity{Data: val})
	return reply.NewIntReply(rowAffected)
}

// ExecGetSet getset k1 v1
func ExecGetSet(db *DB, args [][]byte) respinterface.Reply {
	key := string(args[0])
	val := args[1]
	entity, exists := db.GetEntity(key)
	db.SetEntity(key, &dict.DataEntity{Data: val})
	if !exists {
		return reply.NewEmptyBulkReply()
	}
	return reply.NewBulkReply(entity.Data.([]byte))
}

// ExecSrtLen strlen key
func ExecSrtLen(db *DB, args [][]byte) respinterface.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.NewEmptyBulkReply()
	}
	val := entity.Data.([]byte)
	return reply.NewIntReply(len(val))
}

func init() {
	RegisterCommand("ping", Ping, 1)
	RegisterCommand("del", ExecDel, -2)
	RegisterCommand("exists", ExecExists, -2)
	RegisterCommand("flushdb", ExecFlushDB, -1)
	RegisterCommand("type", ExecType, 2)
	RegisterCommand("rename", ExecRename, 3)
	RegisterCommand("renamenx", ExecRenameNx, 3)
	RegisterCommand("keys", ExecKeys, 2)
	RegisterCommand("get", ExecGet, 2)
	RegisterCommand("set", ExecSet, 3)
	RegisterCommand("setnx", ExecSetNx, 3)
	RegisterCommand("getset", ExecGetSet, 3)
	RegisterCommand("strlen", ExecSrtLen, 2)
}
