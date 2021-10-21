package SealP2P

import (
	"github.com/SealSC/SealP2P/conn"
	"github.com/SealSC/SealP2P/conn/msg"
	"crypto/rsa"
)

type Discoverer interface {
	Listen() error
	Stop()
	Online(ip []string) error
	Offline() error
	SendMsg(payload *msg.Payload) error
	On(func(req *msg.Payload) *msg.Payload)
}

type ConnedNode struct {
	NodeID string
	PubKey *rsa.PublicKey
	pk     *rsa.PrivateKey
	conn   conn.Connect
	IP     []string
	connIP string
}

type Connector interface {
	Listen() error
	Stop()
	NodeList() (list []ConnedNode)
	GetConn(key string) (conn.Connect, bool)
	CloseAndDel(key string)
	DoConn(nodeID string, port int, ip []string) error
	On(func(req *msg.Payload) *msg.Payload)
}
