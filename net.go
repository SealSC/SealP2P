package SealP2P

import (
	"net"
	"github.com/SealSC/SealP2P/conn"
	"github.com/SealSC/SealP2P/conn/msg"
	"github.com/SealSC/SealP2P/conf"
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

func NewNetwork(conf conf.Config, h Handler) (*Network, error) {
	t, err := NewTcpService(conf)
	if err != nil {
		return nil, err
	}
	t.On(h.doHandle)
	m, err := NewMulticast(conf)
	if err != nil {
		return nil, err
	}
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
