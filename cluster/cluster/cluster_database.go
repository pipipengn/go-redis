package cluster

import (
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
	"go-redis/cluster/client"
	"go-redis/cluster/consistenthash"
	"go-redis/config"
	"go-redis/core/database"
	idatabase "go-redis/core/database/interface"
	respinterface "go-redis/resp/interface"
	"go-redis/resp/reply"
	"go-redis/utils/cmdconv"
	"strconv"
)

type Database struct {
	self            string
	nodes           []string
	peerPicker      *consistenthash.NodeMap
	nodeConnections map[string]*pool.ObjectPool
	databaseCore    idatabase.Interface
}

func NewDatabase() *Database {
	cluster := &Database{
		self:            config.Config.Self,
		peerPicker:      consistenthash.NewEmptyNodeMap(nil),
		nodeConnections: map[string]*pool.ObjectPool{},
		databaseCore:    database.New(),
	}

	nodes := append(config.Config.Peers, config.Config.Self)
	cluster.peerPicker.AddNodes(nodes...)
	cluster.nodes = nodes

	for _, peer := range config.Config.Peers {
		objectPool := pool.NewObjectPoolWithDefaultConfig(context.Background(), client.NewConnectionFactory(peer))
		cluster.nodeConnections[peer] = objectPool
	}
	return cluster
}

func (d *Database) Exec(client respinterface.Connection, args [][]byte) respinterface.Reply {
	//TODO implement me
	panic("implement me")
}

func (d *Database) Close() {
	//TODO implement me
	panic("implement me")
}

func (d *Database) AfterClientClose(client respinterface.Connection) {
	//TODO implement me
	panic("implement me")
}

// Communication ================================================================================

func (d *Database) getPeerClient(peer string) (*client.Client, error) {
	objectPool, ok := d.nodeConnections[peer]
	if !ok {
		return nil, errors.New("get peer client: connot find peer connection")
	}

	object, err := objectPool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}

	c, ok := object.(*client.Client)
	if !ok {
		return nil, errors.New("get peer client: type error")
	}

	return c, nil
}

func (d *Database) returnPeerClient(peer string, client *client.Client) error {
	objectPool, ok := d.nodeConnections[peer]
	if !ok {
		return errors.New("return peer client: connot find peer connection")
	}

	return objectPool.ReturnObject(context.Background(), client)
}

func (d *Database) relay(peer string, connection respinterface.Connection, cmd [][]byte) respinterface.Reply {
	if peer == d.self {
		return d.databaseCore.Exec(connection, cmd)
	}

	peerClient, err := d.getPeerClient(peer)
	if err != nil {
		return reply.NewErrReply(err.Error())
	}
	defer func() {
		_ = d.returnPeerClient(peer, peerClient)
	}()

	peerClient.Send(cmdconv.ToCmdLineStrings("SELECT", strconv.Itoa(connection.GetDBIndex())))
	return peerClient.Send(cmd)
}

func (d *Database) broadcast(connection respinterface.Connection, cmd [][]byte) map[string]respinterface.Reply {
	m := map[string]respinterface.Reply{}
	for _, node := range d.nodes {
		m[node] = d.relay(node, connection, cmd)
	}
	return m
}
