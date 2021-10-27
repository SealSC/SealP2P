package SealP2P

import (
	"net"
	"github.com/SealSC/SealP2P/conn"
	"github.com/SealSC/SealP2P/conn/msg"
	"errors"
)

type OnlineInfo struct {
	NodeID  string
	IP      []string
	Port    int
	Version string
}

type Network struct {
	Discoverer
	Connector
}

func NewNetwork(nodeID string, h Handler) (*Network, error) {
	t, err := NewTcpService(nodeID)
	if err != nil {
		return nil, err
	}
	t.On(h.doHandle)
	m := NewMulticast()
	m.On(h.doHandle)
	return &Network{Discoverer: m, Connector: t}, err
}

func ListenMulticastUDP(network string, ifi *net.Interface, gaddr *net.UDPAddr) (conn.UDPConnect, error) {
	udp, err := net.ListenMulticastUDP(network, ifi, gaddr)
	if err != nil {
		return nil, err
	}
	return conn.NewUDPConnect(udp, true), err
}

func SendUdp(address string, p *msg.Message) error {
	if p == nil {
		return nil
	}
	dial, err := net.Dial("udp", address)
	if err != nil {
		return err
	}
	defer dial.Close()
	connect := conn.NewUDPConnect(dial, false)
	connect.Write(p)
	return nil
}

type TCPListener struct {
	l      net.Listener
	nodeID string
}

func ListenTCP(nodeID, address string) (*TCPListener, error) {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	if nodeID == "" {
		return nil, errors.New("nodeID empty")
	}
	listener := &TCPListener{l: listen, nodeID: nodeID}
	return listener, nil
}

func (l *TCPListener) accept() (conn.TCPConnect, error) {
	accept, err := l.l.Accept()
	if err != nil {
		return nil, err
	}
	return conn.NewTCPConnect(accept, false, l.nodeID), err
}

func (l *TCPListener) Close() error {
	return l.l.Close()
}
