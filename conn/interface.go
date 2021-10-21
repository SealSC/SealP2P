package conn

import (
	"github.com/SealSC/SealP2P/conn/msg"
)

type Connect interface {
	Close() error
	Closed() bool
	Write(payload *msg.Payload)
	Read() *msg.Payload
	writeByte([]byte) error
	readByte() ([]byte, error)
}
