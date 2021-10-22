package SealP2P

import (
	"fmt"
	"sync"
	"log"
	"github.com/SealSC/SealP2P/conn"
	"github.com/SealSC/SealP2P/conn/msg"
	"net"
	"errors"
)

var (
	tcpPort = 3333
)

type TcpService struct {
	nodeID  string
	cache   map[string]*ConnedNode
	lock    sync.Mutex
	started bool
	f       func(req *msg.Payload) *msg.Payload
}

func (t *TcpService) Started() bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.started
}

func (t *TcpService) On(f func(req *msg.Payload) *msg.Payload) {
	t.f = f
}

func NewTcpService(nodeID string) (*TcpService, error) {
	if nodeID == "" {
		return nil, errors.New("newTcpService nodeID empty")
	}
	return &TcpService{nodeID: nodeID, cache: map[string]*ConnedNode{}}, nil
}

func (t *TcpService) Stop() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.started = false
}

func (t *TcpService) Listen() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.started {
		return nil
	}
	t.started = true
	listener, err := ListenTCP(t.nodeID, fmt.Sprintf(":%d", tcpPort))
	if err != nil {
		return err
	}
	go func() {
		for t.started {
			conn, err := listener.accept()
			if err != nil {
				panic(fmt.Errorf("listener accept err:%v", err))
			}
			if err := t.goConn(conn); err != nil {
				continue
			}
		}
		_ = listener.Close()
	}()
	return nil
}

func (t *TcpService) goConn(conn conn.TCPConnect) error {
	if err := conn.Handshake(); err != nil {
		conn.Close()
		return err
	}
	t.saveConn(&ConnedNode{
		NodeID: conn.RemoteNodeID(),
		conn:   conn,
		Addr:   conn.RemoteAddr(),
	})
	go func() {
		for t.started && !conn.Status().Closed() {
			req := conn.Read()
			if req != nil && t.f != nil {
				resp := t.f(req)
				conn.Write(resp)
			}
		}
	}()
	return nil
}

func (t *TcpService) saveConn(info *ConnedNode) {
	t.lock.Lock()
	t.lock.Unlock()
	t.cache[info.NodeID] = info
}
func (t *TcpService) NodeList() (list []ConnedNode) {
	t.lock.Lock()
	t.lock.Unlock()
	for _, s := range t.cache {
		list = append(list, *s)
	}
	return list
}

func (t *TcpService) GetConn(key string) (conn.Connect, bool) {
	t.lock.Lock()
	t.lock.Unlock()
	info, ok := t.cache[key]
	if info == nil && ok {
		delete(t.cache, key)
		return nil, false
	}
	if !ok {
		return nil, false
	}
	if info.conn != nil && info.conn.Status().Closed() {
		info.conn.Close()
		delete(t.cache, key)
		return nil, false
	}
	return info.conn, ok
}

func (t *TcpService) CloseAndDel(key string) {
	t.lock.Lock()
	t.lock.Unlock()
	if t.cache == nil {
		t.cache = map[string]*ConnedNode{}
		return
	}
	info := t.cache[key]
	if info.conn != nil {
		info.conn.Close()
	}
	delete(t.cache, key)
}

func (t *TcpService) DialTCP(addr string) (conn.TCPConnect, error) {
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn.NewTCPConnect(dial, true, t.nodeID), nil
}

func (t *TcpService) DoConn(nodeID string, port int, ip []string) error {
	var (
		con conn.TCPConnect
		err error
	)
	if _, ok := t.GetConn(nodeID); ok {
		return nil
	}
	for i := range ip {
		addr := fmt.Sprintf("%s:%d", ip[i], port)
		if con, err = t.DialTCP(addr); err != nil {
			log.Printf("dial node err:%s addr:%v err:%v", nodeID, addr, err)
			continue
		}
		if con != nil {
			break
		}
	}
	if con == nil {
		return fmt.Errorf("cannot connect node:%s ips:%v port:%d", nodeID, ip, port)
	}
	return t.goConn(con)
}
