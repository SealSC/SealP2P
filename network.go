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
	OnMessage(p *msg.Message) *msg.Message
}

type NetNode interface {
	SetMessenger(Messenger)
	GetPubKey() []byte
	GetNodeStatus() NodeStatus
	GetNodeID() string
	GetNodeList() []NodeInfo
	Join() error
	Leave() error
	SendMsg(data *msg.Message) error
	MulticastMsg(data *msg.Message)
	BroadcastMsg(data *msg.Message) error
	MsgProcessorRegister(string, func(req *msg.Message) *msg.Message)
}

type Handler interface {
	SetMessenger(Messenger)
	RegisterHandler(string, func(req *msg.Message) *msg.Message)
	doHandle(req *msg.Message) *msg.Message
}
