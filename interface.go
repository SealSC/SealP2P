package SealP2P

import (
	"github.com/SealSC/SealP2P/conn"
	"github.com/SealSC/SealP2P/conn/msg"
	"crypto/rsa"
)

type Discoverer interface {
	Listener
	Online(ip []string) error
	Offline() error
	SendMsg(payload *msg.Payload) error
	On(func(req *msg.Payload) *msg.Payload)
}

type ConnedNode struct {
	NodeID string
	PubKey *rsa.PublicKey
	pk     *rsa.PrivateKey
	conn   conn.TCPConnect
	Addr   string
}

type Connector interface {
	Listener
	NodeList() (list []ConnedNode)
	GetConn(key string) (conn.Connect, bool)
	CloseAndDel(key string)
	DoConn(nodeID string, port int, ip []string) error
	On(func(req *msg.Payload) *msg.Payload)
}
type Listener interface {
	Listen() error
	Started() bool
	Stop()
}
