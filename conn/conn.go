package conn

import (
	"net"
	"errors"
)

func NewTCPConnect(c net.Conn, client bool, lNode string) (TCPConnect, error) {
	if lNode == "" {
		return nil, errors.New("lNode is nil")
	}
	t := map[bool]Type{false: TypeService, true: TypeClient}[client]
	return &DefaultTCPConnect{c: c, t: t, stat: StatusInit, lNode: lNode}, nil
}
func NewUDPConnect(c net.Conn, multicast bool) UDPConnect {
	return &DefaultUDPConnect{c: c, multicast: multicast}
}
