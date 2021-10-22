package conn

import (
	"github.com/SealSC/SealP2P/conn/msg"
)

type Type string

const (
	TypeClient   Type = "client"
	TypeService  Type = "service"
	TypeLoopback Type = "Loopback"
)

type Status int

func (s Status) Init() bool {
	return s == StatusInit
}
func (s Status) Active() bool {
	return s == StatusActive
}
func (s Status) Closed() bool {
	return s == StatusClosed
}

const (
	StatusInit      Status = 0
	StatusHandshake Status = 1
	StatusActive    Status = 2
	StatusClosed    Status = 3
)

type UDPConnect interface {
	Multicast() bool
	Connect
}

type TCPConnect interface {
	RemoteNodeID() string
	RemoteAddr() string
	LocalNodeID() string
	Handshake() error
	Type() Type
	Status() Status
	Connect
}

type Connect interface {
	Write(payload *msg.Payload)
	Read() *msg.Payload
	Close() error
}
