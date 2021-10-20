package SealP2P

import (
	"net"
	"github.com/SealSC/SealP2P/conn"
	"github.com/SealSC/SealP2P/conn/msg"
)

type OnlineInfo struct {
	NodeID  string
	IP      []string
	Port    int
	Version string
}

type Network struct {
	//服务发现
	Discoverer
	//服务连接
	Connector
	//服务处理
}

func NewNetwork(h Handler) *Network {
	t := NewTcpService()
	t.On(h.doHandle)
	m := NewMulticast()
	m.On(h.doHandle)
	return &Network{Discoverer: m, Connector: t}
}

func Listen(network, address string) (*Listener, error) {
	listen, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return &Listener{listen}, nil
}

func ListenMulticastUDP(network string, ifi *net.Interface, gaddr *net.UDPAddr) (conn.Connect, error) {
	udp, err := net.ListenMulticastUDP(network, ifi, gaddr)
	if err != nil {
		return nil, err
	}
	return conn.NewConnect(udp), err
}

func SendUdp(address string, p *msg.Payload) error {
	if p == nil {
		return nil
	}
	dial, err := net.Dial("udp", address)
	if err != nil {
		return err
	}
	defer dial.Close()
	connect := conn.NewConnect(dial)
	return connect.Write(p)
}

type Listener struct {
	l net.Listener
}

func (l *Listener) accept() (conn.Connect, error) {
	accept, err := l.l.Accept()
	if err != nil {
		return nil, err
	}
	return conn.NewConnect(accept), err
}

func (l *Listener) Close() error {
	return l.l.Close()
}
