package SealP2P

import (
	"encoding/json"
	"github.com/SealSC/SealP2P/conn/msg"
)

func NewJsonPayload(body interface{}) (*msg.Payload, error) {
	marshal, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return &msg.Payload{FromID: localNode.nodeID, Version: version, Body: marshal}, nil
}
func NewPayload() *msg.Payload {
	return &msg.Payload{FromID: localNode.nodeID, Version: version}
}
