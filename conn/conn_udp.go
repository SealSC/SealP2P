package conn

import (
	"github.com/SealSC/SealP2P/conn/msg"
	"github.com/SealSC/SealP2P/tools/varint"
	"bufio"
	"encoding/binary"
	"io/ioutil"
	"io"
	"net"
)

type DefaultUDPConnect struct {
	multicast bool
	c         net.Conn
}

func (d *DefaultUDPConnect) Close() error {
	if d.c != nil {
		return d.c.Close()
	}
	return nil
}

func (d *DefaultUDPConnect) Multicast() bool {
	return d.multicast
}

func (d *DefaultUDPConnect) Write(payload *msg.Message) {
	if payload == nil {
		return
	}
	bytes := payload.PackByte()
	size := varint.New(int64(len(bytes)))
	if d.c != nil {
		d.c.Write(append(size, bytes...))
	}
}

func (d *DefaultUDPConnect) Read() *msg.Message {
	reader := bufio.NewReader(d.c)
	i, err := binary.ReadVarint(reader)
	if err != nil {
		return nil
	}
	all, err := ioutil.ReadAll(io.LimitReader(reader, i))
	if err != nil {
		return nil
	}
	payload := &msg.Message{}
	err = payload.UNPackByte(all)
	if payload.Action == "" {
		return nil
	}
	return payload
}
