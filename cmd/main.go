package main

import (
	"github.com/SealSC/SealP2P"
	"time"
	"fmt"
	"log"
	"net"
)

func main() {
	conn := net.TCPConn{}
	conn.Close()
	log.SetFlags(log.Llongfile | log.Ltime)
	node := SealP2P.LocalNode()
	err := node.Join()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 2)
	payload, err := SealP2P.NewJsonPayload("{}")
	if err != nil {
		panic(err)
	}
	err = node.SendMsg(payload)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
	fmt.Println("local node id:", node.GetNodeID())
	fmt.Println("node list :", node.GetNodeList())
	node.MulticastMsg(payload)
	node.BroadcastMsg(payload)
	time.Sleep(time.Second)
	err = node.Leave()
	if err != nil {
		panic(err)
	}
	select {}
}
