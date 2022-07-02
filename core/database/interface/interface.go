package idatabase

import respinterface "go-redis/resp/interface"

type Face interface {
	Exec(client respinterface.Connection, args [][]byte) respinterface.Reply
	Close()
	AfterClientClose(client respinterface.Connection)
}
