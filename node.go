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
	"github.com/SealSC/SealP2P/conf"
)

type Node struct {
	conf    *conf.Config
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
	return n.conf.ID
}
func (n *Node) GetNodeStatus() NodeStatus {
	n.status.Dis = n.network.Discoverer.Started()
	n.status.Ser = n.network.Connector.Started()
	return n.status
}

func (n *Node) GetNodeList() (list []NodeInfo) {
	nodeList := n.network.NodeList()
	for _, node := range nodeList {
		list = append(list, NodeInfo{
			ID:   node.NodeID,
			Addr: node.Addr,
			Type: node.conn.Type(),
		})
	}
	return list
}

func (n *Node) Join() error {
	if !n.conf.ClientOnly {
		err := n.network.Connector.Listen()
		if err != nil {
			return err
		}
	}

	err := n.network.Discoverer.Listen()
	if err != nil {
		return err
	}
	err = n.network.Online(n.status.IP)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) Leave() error {
	if err := n.network.Discoverer.Offline(); err != nil {
		return err
	}
	n.network.Discoverer.Stop()
	n.network.Connector.Stop()
	return nil
}

func (n *Node) SendMsg(data *msg.Message) error {
	if data == nil {
		return nil
	}
	if len(data.ToID) < 1 {
		return errors.New("invalid send destination")
	}
	data.Action = msg.ActionDail
	conn, ok := n.network.Connector.GetConn(data.ToID[0])
	if ok {
		conn.Write(data)
	}
	return nil
}

func (n *Node) MulticastMsg(data *msg.Message) {
	if data == nil {
		return
	}
	idSet := map[string]struct{}{}
	for i := range data.ToID {
		idSet[data.ToID[i]] = struct{}{}
	}
	data.Action = msg.ActionMulticast
	for s := range idSet {
		conn, ok := n.network.Connector.GetConn(s)
		if ok {
			conn.Write(data)
		}
	}

}

func (n *Node) BroadcastMsg(data *msg.Message) error {
	if data == nil {
		return nil
	}
	data.Action = msg.ActionBroadcast
	return n.network.Discoverer.SendMsg(data)
}

func (n *Node) MsgProcessorRegister(router string, f func(req *msg.Message) *msg.Message) {
	n.h.RegisterHandler(router, f)
}

var localNode *Node
var newLock sync.Mutex

func LocalNode() NetNode {
	return localNode
}

func InitLocalNode(conf *conf.Config) error {
	newLock.Lock()
	newLock.Unlock()
	if localNode != nil {
		return nil
	}
	if conf == nil {
		return errors.New("conf is nil")
	}
	n := &Node{h: NewDefaultHandler(), conf: conf}
	key, err := readRSA(conf.PKFile)
	if err != nil {
		return err
	}
	n.key = key
	available, err := ip.Available()
	if err != nil {
		return err
	}
	conf.ID = grsa.PubSha1(key)
	n.status = NodeStatus{
		ID: conf.ID,
		IP: available,
	}
	network, err := NewNetwork(conf, n.h)
	if err != nil {
		return err
	}
	n.network = network
	localNode = n
	return nil
}

func readRSA(pkFile string) (pk *rsa.PrivateKey, err error) {
	//Read the configuration if it exists
	if gio.FileExist(pkFile) {
		return grsa.LoadFile(pkFile)
	}
	pk, err = grsa.RandKey()
	if err != nil {
		return nil, err
	}
	err = grsa.SaveFile(pkFile, pk)
	return pk, err
}
