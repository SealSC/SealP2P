package SealP2P

import (
	"fmt"
	"sync"
	"log"
	"github.com/SealSC/SealP2P/conn"
	"github.com/SealSC/SealP2P/conn/msg"
	"crypto/rsa"
	"net"
)

var (
	tcpPort = 3333
)

type ConnedNode struct {
	NodeID string
	PubKey *rsa.PublicKey
	pk     *rsa.PrivateKey
	conn   conn.Connect
	connIP string
}

type TcpService struct {
	cache   map[string]ConnedNode
	lock    sync.Mutex
	started bool
	f       func(req *msg.Payload) *msg.Payload
}

func (t *TcpService) On(f func(req *msg.Payload) *msg.Payload) {
	t.f = f
}

func NewTcpService() *TcpService {
	return &TcpService{cache: map[string]ConnedNode{}}
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
	listener, err := Listen("tcp", fmt.Sprintf(":%d", tcpPort))
	if err != nil {
		return err
	}
	go func() {
		for t.started {
			conn, err := listener.accept()
			if err != nil {
				continue
			}
			go t.onConn(conn)
		}
		_ = listener.Close()
	}()
	return nil
}
func (t *TcpService) onConn(conn conn.Connect) {
	for t.started {
		req, err := conn.Read()
		if err != nil {
			log.Println("conn read payload err:", err)
			continue
		}
		if t.f != nil {
			resp := t.f(req)
			conn.Write(resp)
		}
	}
}

func (t *TcpService) SaveConn(info ConnedNode) {
	t.lock.Lock()
	t.lock.Unlock()
	t.cache[info.NodeID] = info
}
func (t *TcpService) NodeList() (list []string) {
	t.lock.Lock()
	t.lock.Unlock()
	for s := range t.cache {
		list = append(list, s)
	}
	return list
}

func (t *TcpService) GetConn(key string) (conn.Connect, bool) {
	t.lock.Lock()
	t.lock.Unlock()
	info, ok := t.cache[key]
	if info.conn != nil && info.conn.Closed() {
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
		t.cache = map[string]ConnedNode{}
		return
	}
	info := t.cache[key]
	if info.conn != nil {
		info.conn.Close()
	}
	delete(t.cache, key)
}

func (t *TcpService) DialTCP(addr string) (conn.Connect, error) {
	dial, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return conn.NewConnect(dial), nil
}

func (t *TcpService) DoConn(nodeID string, port int, ip []string) error {
	var (
		con conn.Connect
		err error
	)
	if _, ok := t.GetConn(nodeID); ok {
		return nil
	}
	node := ConnedNode{NodeID: nodeID}
	for i := range ip {
		addr := fmt.Sprintf("%s:%d", ip[i], port)
		if con, err = t.DialTCP(addr); err != nil {
			log.Printf("dial node err:%s addr:%v err:%v", nodeID, addr, err)
			continue
		}
		if con != nil {
			node.connIP = ip[i]
			break
		}
	}
	if con == nil {
		return fmt.Errorf("cannot connect node:%s ips:%v port:%d", node.NodeID, ip, port)
	}
	node.conn = con
	t.SaveConn(node)
	go t.onConn(con)
	return nil
}
