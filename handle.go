package SealP2P

import (
	"log"
	"encoding/json"
	"github.com/SealSC/SealP2P/conn/msg"
)

var DefaultHandleMap = map[string]func(payload *msg.Message) *msg.Message{
	msg.Join: func(request *msg.Message) *msg.Message {
		log.Println("join:", request.FromID, request)
		info := OnlineInfo{}
		err := json.Unmarshal(request.Payload, &info)
		if err != nil {
			panic(err)
		}
		if err = localNode.network.DoConn(info.NodeID, info.Port, info.IP); err != nil {
			log.Println("online err", err)
		}
		return nil
	},
	msg.Leave: func(request *msg.Message) *msg.Message {
		log.Println("msg.Leave:", request.FromID)
		localNode.network.CloseAndDel(request.FromID)
		return nil
	},
}

type DefaultHandler struct {
	customMap map[string]func(payload *msg.Message) *msg.Message
	m         Messenger
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{customMap: map[string]func(payload *msg.Message) *msg.Message{}}
}
func (d *DefaultHandler) SetMessenger(m Messenger) {
	d.m = m
}

func (d *DefaultHandler) RegisterHandler(key string, f func(req *msg.Message) *msg.Message) {
	if f == nil {
		return
	}
	if d.customMap == nil {
		d.customMap = map[string]func(payload *msg.Message) *msg.Message{}
	}
	d.customMap[key] = f
}

func (d *DefaultHandler) doHandle(req *msg.Message) *msg.Message {
	if req == nil {
		return nil
	}
	if req.FromID == localNode.GetNodeID() {
		return nil
	}
	switch req.Type {
	case msg.Dail, msg.Multicast, msg.Broadcast:
		if d.m != nil {
			return d.m.OnMessage(req)
		}
		return nil
	default:
		if d.m != nil {
			d.m.OnMessage(req)
		}
	}
	if f := DefaultHandleMap[req.Type]; f != nil {
		return f(req)
	}
	if d.customMap == nil {
		return nil
	}
	if f := d.customMap[req.Type]; f != nil {
		return f(req)
	}
	return nil
}
