package conn

import (
	"github.com/SealSC/SealP2P/conn/msg"
)

type Connect interface {
	Close() error
	Closed() bool
	Write(payload *msg.Payload) error
	Read() (*msg.Payload, error)
	writeByte([]byte) error
	readByte() ([]byte, error)
}
