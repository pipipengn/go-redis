package client

import (
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
)

type ConnectionFactory struct {
	PeerAddr string
}

func NewConnectionFactory(addr string) *ConnectionFactory {
	return &ConnectionFactory{PeerAddr: addr}
}

func (f *ConnectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	client, err := NewClient(f.PeerAddr)
	if err != nil {
		return nil, err
	}
	client.Start()
	return pool.NewPooledObject(client), nil
}

func (f *ConnectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	client, ok := object.Object.(Client)
	if !ok {
		return errors.New("connection pool type error")
	}
	client.Close()
	return nil
}

func (f *ConnectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (f *ConnectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (f *ConnectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
