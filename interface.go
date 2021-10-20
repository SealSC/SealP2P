package SealP2P

import (
	"github.com/SealSC/SealP2P/conn"
	"github.com/SealSC/SealP2P/conn/msg"
)

type Discoverer interface {
	Listen() error
	Stop()
	Online() error
	Offline() error
	SendMsg(payload *msg.Payload) error
	On(func(req *msg.Payload) *msg.Payload)
}

type Connector interface {
	Listen() error
	Stop()
	NodeList() (list []string)
	GetConn(key string) (conn.Connect, bool)
	CloseAndDel(key string)
	DoConn(nodeID string, port int, ip []string) error
	On(func(req *msg.Payload) *msg.Payload)
}
