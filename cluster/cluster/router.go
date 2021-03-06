package cluster

import (
	respinterface "go-redis/resp/interface"
	"go-redis/resp/reply"
	"go-redis/utils/cmdconv"
	"strings"
)

type cmdRelayFunc func(cluster *Database, connection respinterface.Connection, cmd [][]byte) respinterface.Reply

var router = make(map[string]cmdRelayFunc)

func RegisterRouter(name string, executor cmdRelayFunc) {
	name = strings.ToLower(name)
	router[name] = executor
}

func defaultFunc(cluster *Database, connection respinterface.Connection, cmd [][]byte) respinterface.Reply {
	peer := cluster.peerPicker.PickNode(string(cmd[1]))
	return cluster.relay(peer, connection, cmd)
}

func localFunc(cluster *Database, connection respinterface.Connection, cmd [][]byte) respinterface.Reply {
	return cluster.databaseCore.Exec(connection, cmd)
}

func renameFunc(cluster *Database, connection respinterface.Connection, cmd [][]byte) respinterface.Reply {
	if len(cmd) != 3 {
		return reply.NewArgNumErrReply("rename")
	}

	key1 := string(cmd[1])
	key2 := string(cmd[2])
	oldPeer := cluster.peerPicker.PickNode(key1)
	newPeer := cluster.peerPicker.PickNode(key2)

	if oldPeer == newPeer {
		return cluster.relay(oldPeer, connection, cmd)
	}

	// oldpeer get val
	line := cmdconv.ToCmdLineStrings("get", key1)
	relayReply := cluster.relay(oldPeer, connection, line)
	if _, ok := relayReply.(*reply.EmptyBulkReply); ok {
		return reply.NewErrReply("key does not exists")
	}

	bulkReply, ok := relayReply.(*reply.BulkReply)
	if !ok {
		return reply.NewErrReply("type error")
	}
	val := string(bulkReply.Arg)

	// newpeer  set key2 val
	line = cmdconv.ToCmdLineStrings("set", key2, val)
	okReply := cluster.relay(newPeer, connection, line)

	// oldpeer  del key1
	line = cmdconv.ToCmdLineStrings("del", key1)
	cluster.relay(oldPeer, connection, line)

	return okReply
}

func renameNxFunc(cluster *Database, connection respinterface.Connection, cmd [][]byte) respinterface.Reply {
	if len(cmd) != 3 {
		return reply.NewArgNumErrReply("rename")
	}

	key1 := string(cmd[1])
	key2 := string(cmd[2])
	oldPeer := cluster.peerPicker.PickNode(key1)
	newPeer := cluster.peerPicker.PickNode(key2)

	if oldPeer == newPeer {
		return cluster.relay(oldPeer, connection, cmd)
	}

	// oldpeer get val
	line := cmdconv.ToCmdLineStrings("get", key1)
	relayReply := cluster.relay(oldPeer, connection, line)
	if _, ok := relayReply.(*reply.EmptyBulkReply); ok {
		return reply.NewErrReply("key does not exists")
	}

	bulkReply, ok := relayReply.(*reply.BulkReply)
	if !ok {
		return reply.NewErrReply("type error")
	}
	val := string(bulkReply.Arg)

	// newpeer  setnx key2 val
	line = cmdconv.ToCmdLineStrings("setnx", key2, val)
	okReply := cluster.relay(newPeer, connection, line)

	// oldpeer  del key1
	line = cmdconv.ToCmdLineStrings("del", key1)
	cluster.relay(oldPeer, connection, line)

	return okReply
}

func flushDbFunc(cluster *Database, connection respinterface.Connection, cmd [][]byte) respinterface.Reply {
	replyMap := cluster.broadcast(connection, cmd)
	for _, execReply := range replyMap {
		if reply.IsErrorReply(execReply) {
			return reply.NewErrReply("connot flush db: " + string(execReply.ToBytes()))
		}
	}
	return reply.NewOkReply()
}

func deleteFunc(cluster *Database, connection respinterface.Connection, cmd [][]byte) respinterface.Reply {
	replyMap := cluster.broadcast(connection, cmd)
	deletedNum := 0
	for _, execReply := range replyMap {
		intReply := execReply.(*reply.IntReply)
		deletedNum += intReply.Code
	}
	return reply.NewIntReply(deletedNum)
}

func init() {
	RegisterRouter("exists", defaultFunc)
	RegisterRouter("type", defaultFunc)
	RegisterRouter("get", defaultFunc)
	RegisterRouter("set", defaultFunc)
	RegisterRouter("setnx", defaultFunc)
	RegisterRouter("getset", defaultFunc)
	RegisterRouter("ping", localFunc)
	RegisterRouter("select", localFunc)
	RegisterRouter("rename", renameFunc)
	RegisterRouter("renamenx", renameNxFunc)
	RegisterRouter("flushdb", flushDbFunc)
	RegisterRouter("delete", deleteFunc)
}
