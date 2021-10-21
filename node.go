package SealP2P

import (
	"crypto/rsa"
	"sync"
	"github.com/SealSC/SealP2P/tools/grsa"
	"crypto/x509"
	"github.com/SealSC/SealP2P/conn/msg"
	"github.com/SealSC/SealP2P/tools/gio"
	"errors"
	"github.com/SealSC/SealP2P/tools/ip"
)

type Node struct {
	key     *rsa.PrivateKey
	network *Network
	h       Handler
	status  NodeStatus
}

func (n *Node) SetMessenger(m Messenger) {
	if n.h != nil {
		n.h.SetMessenger(m)
	}
}
func (n *Node) GetPubKey() []byte {
	return x509.MarshalPKCS1PublicKey(&n.key.PublicKey)
}

func (n *Node) GetNodeID() string {
	return n.status.ID
}
func (n *Node) GetNodeStatus() NodeStatus {
	return n.status
}

func (n *Node) GetNodeList() (list []NodeInfo) {
	nodeList := n.network.NodeList()
	for _, node := range nodeList {
		list = append(list, NodeInfo{
			ID:     node.NodeID,
			IP:     node.IP,
			ConnIP: node.connIP,
		})
	}
	return list
}

func (n *Node) Join() error {
	err := n.network.Connector.Listen()
	if err != nil {
		return err
	}

	err = n.network.Discoverer.Listen()
	if err != nil {
		return err
	}
	err = n.network.Online(n.status.IP)
	if err != nil {
		return err
	}
	n.status.Online = true
	return nil
}

func (n *Node) Leave() error {
	if err := n.network.Discoverer.Offline(); err != nil {
		return err
	}
	n.network.Discoverer.Stop()
	n.network.Connector.Stop()
	n.status.Online = false
	return nil
}

func (n *Node) SendMsg(data *msg.Payload) error {
	if data == nil {
		return nil
	}
	if len(data.ToID) < 1 {
		return errors.New("invalid send destination")
	}
	data.Path = msg.Dail
	conn, ok := n.network.Connector.GetConn(data.ToID[0])
	if ok {
		conn.Write(data)
	}
	return nil
}

func (n *Node) MulticastMsg(data *msg.Payload) {
	if data == nil {
		return
	}
	idSet := map[string]struct{}{}
	for i := range data.ToID {
		idSet[data.ToID[i]] = struct{}{}
	}
	data.Path = msg.Multicast
	for s := range idSet {
		conn, ok := n.network.Connector.GetConn(s)
		if ok {
			conn.Write(data)
		}
	}

}

func (n *Node) BroadcastMsg(data *msg.Payload) error {
	if data == nil {
		return nil
	}
	data.Path = msg.Broadcast
	return n.network.Discoverer.SendMsg(data)
}

func (n *Node) MsgProcessorRegister(router string, f func(req *msg.Payload) *msg.Payload) {
	n.h.RegisterHandler(router, f)
}

var localNode *Node
var newLock sync.Mutex

func LocalNode() NetNode {
	return localNode
}

func init() {
	var once = sync.Once{}
	once.Do(func() {
		node, err := newLocalNode("SealP2PPK")
		if err != nil {
			panic("newLocalNode:" + err.Error())
		}
		localNode = node
	})
}

func newLocalNode(pkFile string) (*Node, error) {
	newLock.Lock()
	newLock.Unlock()
	if localNode != nil {
		return localNode, nil
	}
	n := &Node{h: NewDefaultHandler()}
	n.network = NewNetwork(n.h)
	key, err := readRSA(pkFile)
	if err != nil {
		return nil, err
	}
	n.key = key
	available, err := ip.Available()
	if err != nil {
		return nil, err
	}
	n.status = NodeStatus{
		ID:     grsa.PubSha1(key),
		IP:     available,
		Online: false,
	}
	return n, nil
}

func readRSA(pkFile string) (pk *rsa.PrivateKey, err error) {
	//Read the configuration if it exists
	if gio.FileExist(pkFile) {
		pk, err = grsa.LoadFile(pkFile)
		return
	}
	pk, err = grsa.RandKey()
	if err != nil {
		return nil, err
	}
	err = grsa.SaveFile(pkFile, pk)
	return pk, err
}
