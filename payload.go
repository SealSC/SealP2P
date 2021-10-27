package SealP2P

import (
	"encoding/json"
	"github.com/SealSC/SealP2P/conn/msg"
)

func NewJsonMessage(body interface{}) (*msg.Message, error) {
	marshal, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return &msg.Message{FromID: localNode.GetNodeID(), Version: version, Payload: marshal}, nil
}
func EmptyMessage() *msg.Message {
	return &msg.Message{FromID: localNode.GetNodeID(), Version: version}
}

func NewPayload(path string) *msg.Message {
	return &msg.Message{FromID: localNode.GetNodeID(), Version: version, Type: path}
}
