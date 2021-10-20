package SealP2P

import (
	"github.com/SealSC/SealP2P/conn/msg"
)

type NetNode interface {
	GetPubKey() []byte
	//GetNodeID 获取本节点标识
	GetNodeID() string
	//GetNodeList 获取网络中的Node列表
	GetNodeList() []string
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
	RegisterHandler(string, func(req *msg.Payload) *msg.Payload)
	doHandle(req *msg.Payload) *msg.Payload
}
