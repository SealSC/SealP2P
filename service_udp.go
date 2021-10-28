package SealP2P

import (
	"net"
	"fmt"
	"log"
	"sync"
	"github.com/SealSC/SealP2P/conn/msg"
	"time"
	"github.com/SealSC/SealP2P/conf"
	"errors"
)

type Multicast struct {
	conf    *conf.Config
	l       sync.Mutex
	f       func(req *msg.Message) *msg.Message
	started bool
}

func (m *Multicast) Started() bool {
	m.l.Lock()
	defer m.l.Unlock()
	return m.started
}

func NewMulticast(conf *conf.Config) (*Multicast, error) {
	if conf == nil {
		return nil, errors.New("conf is nil")
	}
	if conf.MulticastPort < 0 || conf.MulticastPort > 65535 {
		return nil, fmt.Errorf("port(%d) err", conf.MulticastPort)
	}
	if ip := net.ParseIP(conf.MulticastAddr); !ip.IsMulticast() {
		return nil, fmt.Errorf("ip(%s) not Multicast", conf.MulticastAddr)
	}
	return &Multicast{conf: conf}, nil
}
func (m *Multicast) On(f func(req *msg.Message) *msg.Message) {
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
	udp, err := ListenMulticastUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(m.conf.MulticastAddr), Port: m.conf.ServerPort})
	if err != nil {
		return err
	}
	go func() {
		for m.started {
			req := udp.Read()
			if req == nil {
				time.Sleep(time.Millisecond * 500)
			}
			go m.doReq(req)
		}
		_ = udp.Close()
	}()

	return nil
}
func (m *Multicast) doReq(p *msg.Message) {
	if p == nil {
		return
	}
	switch p.Action {
	case msg.ActionJoin, msg.ActionLeave, msg.ActionBroadcast:
		if m.f != nil {
			m.f(p)
		}
	default:
		log.Println("ActionMulticast don't know path:", p)
	}
}
func (m *Multicast) Offline() (err error) {
	payload, err := NewJsonMessage(nil)
	if err != nil {
		return err
	}
	payload.Action = msg.ActionLeave
	return m.SendMsg(payload)
}

func (m *Multicast) Online(ip []string) (err error) {
	payload, err := NewJsonMessage(OnlineInfo{
		NodeID:  localNode.GetNodeID(),
		IP:      ip,
		Port:    m.conf.ServerPort,
		Version: version,
	})
	if err != nil {
		return err
	}
	payload.Action = msg.ActionJoin
	return m.SendMsg(payload)
}
func (m *Multicast) SendMsg(p *msg.Message) (err error) {
	if p == nil {
		return nil
	}
	return SendUdp(fmt.Sprintf("%s:%d", m.conf.MulticastAddr, m.conf.MulticastPort), p)
}
