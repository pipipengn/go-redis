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

	// oldpeer取出val
	line := cmdconv.ToCmdLineStrings("get", key1)
	relayReply := cluster.relay(oldPeer, connection, line)
	if reply.IsErrorReply(relayReply) {
		return relayReply
	}

	s := string(relayReply.ToBytes())
	split := strings.Split(s, "\r\n")
	val := split[1]

	// newpeer执行 set key2 val
	line = cmdconv.ToCmdLineStrings("set", key2, val)
	ok := cluster.relay(newPeer, connection, line)

	// oldpeer执行 del key1
	line = cmdconv.ToCmdLineStrings("del", key1)
	cluster.relay(oldPeer, connection, line)

	return ok
}

func init() {
	RegisterRouter("exists", defaultFunc)
	RegisterRouter("type", defaultFunc)
	RegisterRouter("get", defaultFunc)
	RegisterRouter("set", defaultFunc)
	RegisterRouter("setnx", defaultFunc)
	RegisterRouter("getset", defaultFunc)
	RegisterRouter("ping", localFunc)
	RegisterRouter("rename", renameFunc)
}
