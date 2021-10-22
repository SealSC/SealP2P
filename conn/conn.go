package conn

import "net"

func NewTCPConnect(c net.Conn, client bool, lNode string) TCPConnect {
	t := map[bool]Type{false: TypeService, true: TypeClient}[client]
	return &DefaultTCPConnect{c: c, t: t, stat: StatusInit, lNode: lNode}
}
func NewUDPConnect(c net.Conn, multicast bool) UDPConnect {
	return &DefaultUDPConnect{c: c, multicast: multicast}
}
