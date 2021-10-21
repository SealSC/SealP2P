package SealP2P

import (
	"log"
	"encoding/json"
	"github.com/SealSC/SealP2P/conn/msg"
)

var DefaultHandleMap = map[string]func(payload *msg.Payload) *msg.Payload{
	msg.PING: func(request *msg.Payload) *msg.Payload {
		newPayload := NewPayload()
		newPayload.Path = msg.PONG
		return newPayload
	},
	msg.PONG: func(request *msg.Payload) *msg.Payload {
		return nil
	},
	msg.Join: func(request *msg.Payload) *msg.Payload {
		log.Println("join:", request.FromID, request)
		info := OnlineInfo{}
		err := json.Unmarshal(request.Body, &info)
		if err != nil {
			panic(err)
		}
		if err = localNode.network.DoConn(info.NodeID, info.Port, info.IP); err != nil {
			log.Println("online err", err)
		}
		return nil
	},
	msg.Leave: func(request *msg.Payload) *msg.Payload { //call local
		log.Println("msg.Leave:", request.FromID)
		localNode.network.CloseAndDel(request.FromID)
		return nil
	},
}

type DefaultHandler struct {
	customMap map[string]func(payload *msg.Payload) *msg.Payload
	m         Messenger
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{customMap: map[string]func(payload *msg.Payload) *msg.Payload{}}
}
func (d *DefaultHandler) SetMessenger(m Messenger) {
	d.m = m
}

func (d *DefaultHandler) RegisterHandler(key string, f func(req *msg.Payload) *msg.Payload) {
	d.customMap[key] = f
}

func (d *DefaultHandler) doHandle(req *msg.Payload) *msg.Payload {
	if req == nil {
		return nil
	}
	//if req.FromID == localNode.GetNodeID() {
	//	return nil
	//}
	switch req.Path {
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
	if f := DefaultHandleMap[req.Path]; f != nil {
		return f(req)
	}
	if f := d.customMap[req.Path]; f != nil {
		return f(req)
	}
	return nil
}
