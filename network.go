package SealP2P

import (
	"github.com/SealSC/SealP2P/conn/msg"
	"github.com/SealSC/SealP2P/conn"
)

type NodeInfo struct {
	ID   string    `json:"id"`
	Addr string    `json:"addr"`
	Type conn.Type `json:"type"`
}

type NodeStatus struct {
	ID  string   `json:"id"`
	IP  []string `json:"ip"`
	Dis bool     `json:"dis"`
	Ser bool     `json:"ser"`
}

type Messenger interface {
	OnMessage(p *msg.Payload) *msg.Payload
}

type NetNode interface {
	SetMessenger(Messenger)
	GetPubKey() []byte
	GetNodeStatus() NodeStatus
	GetNodeID() string
	GetNodeList() []NodeInfo
	Join() error
	Leave() error
	SendMsg(data *msg.Payload) error
	MulticastMsg(data *msg.Payload)
	BroadcastMsg(data *msg.Payload) error
	MsgProcessorRegister(string, func(req *msg.Payload) *msg.Payload)
}

type Handler interface {
	SetMessenger(Messenger)
	RegisterHandler(string, func(req *msg.Payload) *msg.Payload)
	doHandle(req *msg.Payload) *msg.Payload
}
