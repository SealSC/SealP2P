package SealP2P

import (
	"github.com/SealSC/SealP2P/conn/msg"
)

type NodeInfo struct {
	ID     string   `json:"id"`
	IP     []string `json:"ip"`
	ConnIP string   `json:"connIP"`
}

type NodeStatus struct {
	ID     string   `json:"id"`
	IP     []string `json:"ip"`
	Online bool     `json:"online"`
}

type Messenger interface {
	OnMessage(p *msg.Payload) *msg.Payload
}

type NetNode interface {
	SetMessenger(Messenger)
	GetPubKey() []byte
	//GetNodeStatus 获取本节点
	GetNodeStatus() NodeStatus
	//GetNodeID 获取本节点标识
	GetNodeID() string
	//GetNodeList 获取网络中的Node列表
	GetNodeList() []NodeInfo
	//Join 加入网络
	Join() error
	//Leave 离开网络
	Leave() error //离开网络
	//SendMsg 单播消息给某个节点
	SendMsg(data *msg.Payload) error
	//MulticastMsg 组播消息给一组节点
	MulticastMsg(data *msg.Payload)
	//BroadcastMsg 全网广播消息
	BroadcastMsg(data *msg.Payload) error
	//MsgProcessorRegister 注册消息处理回调函数
	MsgProcessorRegister(string, func(req *msg.Payload) *msg.Payload)
}

type Handler interface {
	SetMessenger(Messenger)
	RegisterHandler(string, func(req *msg.Payload) *msg.Payload)
	doHandle(req *msg.Payload) *msg.Payload
}
