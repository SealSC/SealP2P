package SealP2P

import (
	"net"
	"fmt"
	"log"
	"sync"
	"github.com/SealSC/SealP2P/conn/msg"
)

var (
	MulticastIP   = "224.0.0.3"
	MulticastPort = 5678
)

type Multicast struct {
	l       sync.Mutex
	f       func(req *msg.Payload) *msg.Payload
	started bool
}

func (m *Multicast) Started() bool {
	m.l.Lock()
	defer m.l.Unlock()
	return m.started
}

func NewMulticast() *Multicast {
	return &Multicast{}
}
func (m *Multicast) On(f func(req *msg.Payload) *msg.Payload) {
	m.f = f
}
func (m *Multicast) Stop() {
	m.l.Lock()
	defer m.l.Unlock()
	m.started = false
}
func (m *Multicast) Listen() error {
	m.l.Lock()
	defer m.l.Unlock()
	if m.started {
		return nil
	}
	m.started = true
	udp, err := ListenMulticastUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(MulticastIP), Port: MulticastPort})
	if err != nil {
		return err
	}
	go func() {
		for m.started {
			req := udp.Read()
			go m.doReq(req)
		}
		_ = udp.Close()
	}()

	return nil
}
func (m *Multicast) doReq(p *msg.Payload) {
	if p == nil {
		return
	}
	switch p.Path {
	case msg.Join, msg.Leave, msg.Broadcast:
		if m.f != nil {
			m.f(p)
		}
	default:
		log.Println("Multicast don't know path:", p)
	}
}
func (m *Multicast) Offline() (err error) {
	payload, err := NewJsonPayload(nil)
	if err != nil {
		return err
	}
	payload.Path = msg.Leave
	return m.SendMsg(payload)
}

func (m *Multicast) Online(ip []string) (err error) {
	newPayload := NewPayload()
	payload, err := NewJsonPayload(OnlineInfo{
		NodeID:  newPayload.FromID,
		IP:      ip,
		Port:    tcpPort,
		Version: newPayload.Version,
	})
	if err != nil {
		return err
	}
	payload.Path = msg.Join
	return m.SendMsg(payload)
}
func (m *Multicast) SendMsg(p *msg.Payload) (err error) {
	if p == nil {
		return nil
	}
	return SendUdp(fmt.Sprintf("%s:%d", MulticastIP, MulticastPort), p)
}
