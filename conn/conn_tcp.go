package conn

import (
	"io"
	"net"
	"bufio"
	"encoding/binary"
	"io/ioutil"
	"sync"
	"github.com/SealSC/SealP2P/tools/varint"
	"github.com/SealSC/SealP2P/conn/msg"
	"errors"
	"time"
	"fmt"
)

type DefaultTCPConnect struct {
	c     net.Conn
	l     sync.Mutex
	stat  Status
	t     Type
	rNode string
	lNode string
}

func (d *DefaultTCPConnect) RemoteAddr() string {
	return d.c.RemoteAddr().String()
}

func (d *DefaultTCPConnect) Handshake() error {
	if !d.stat.Init() {
		return errors.New("conn uninitialized")
	}
	d.l.Lock()
	defer d.l.Unlock()
	payload := &msg.Payload{
		Version: "-",
		FromID:  d.lNode,
		Path:    msg.PING,
	}
	d.Write(payload)
	var c = make(chan *msg.Payload, 1)
	select {
	case <-time.After(time.Second * 10):
		return errors.New("handshake read time out")
	case c <- d.Read():
	}
	read := <-c
	if read.Path != msg.PING {
		return fmt.Errorf("handshake response path is \"%s\"", read.Path)
	}
	d.rNode = read.FromID
	if read.FromID == d.lNode {
		d.t = TypeLoopback
	} else {
		d.stat = StatusActive
	}
	return nil
}

func (d *DefaultTCPConnect) RemoteNodeID() string {
	return d.rNode
}

func (d *DefaultTCPConnect) LocalNodeID() string {
	return d.lNode
}

func (d *DefaultTCPConnect) Type() Type {
	return d.t
}

func (d *DefaultTCPConnect) death(err error) bool {
	if err == nil {
		return false
	}
	if opError, ok := err.(*net.OpError); ok {
		err = opError.Err
	}
	if err == net.ErrClosed ||
		err == io.ErrUnexpectedEOF ||
		err == io.EOF {
		return true
	}
	return false
}

func (d *DefaultTCPConnect) Write(payload *msg.Payload) {
	if payload == nil {
		return
	}
	err := d.writeByte(payload.PackByte())
	if d.death(err) {
		d.Close()
	}
}

func (d *DefaultTCPConnect) Read() *msg.Payload {
	read, err := d.doRead()
	if d.death(err) {
		d.Close()
		return nil
	}
	return read
}

func (d *DefaultTCPConnect) doRead() (*msg.Payload, error) {
	pkg, err := d.readByte()
	if err != nil {
		return nil, err
	}
	payload := &msg.Payload{}
	err = payload.UNPackByte(pkg)
	if err != nil {
		return nil, err
	}
	if payload.Path == "" {
		return nil, nil
	}
	return payload, err
}

func (d *DefaultTCPConnect) writeByte(bytes []byte) error {
	size := varint.New(int64(len(bytes)))
	_, err := d.c.Write(append(size, bytes...))
	return err
}

func (d *DefaultTCPConnect) readByte() ([]byte, error) {
	reader := bufio.NewReader(d.c)
	i, err := binary.ReadVarint(reader)
	if err != nil {
		return nil, err
	}
	all, err := ioutil.ReadAll(io.LimitReader(reader, i))
	return all, err
}

func (d *DefaultTCPConnect) Close() error {
	d.l.Lock()
	defer d.l.Unlock()
	if d.stat.Closed() {
		return nil
	}
	err := d.c.Close()
	if err != nil {
		d.stat = StatusClosed
	}
	return err
}
func (d *DefaultTCPConnect) Status() Status {
	d.l.Lock()
	defer d.l.Unlock()
	return d.stat
}
