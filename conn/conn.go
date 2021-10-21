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
)

type DefaultConnect struct {
	c     net.Conn
	l     sync.Mutex
	close bool
}

func (d *DefaultConnect) death(err error) bool {
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

func (d *DefaultConnect) Write(payload *msg.Payload) {
	if payload == nil {
		return
	}
	err := d.writeByte(payload.PackByte())
	if d.death(err) {
		d.Close()
	}
}

func (d *DefaultConnect) Read() *msg.Payload {
	read, err := d.doRead()
	if d.death(err) {
		d.Close()
		return nil
	}
	return read
}

func (d *DefaultConnect) doRead() (*msg.Payload, error) {
	pkg, err := d.readByte()
	if err != nil {
		return nil, err
	}
	payload := &msg.Payload{}
	err = payload.UNPackByte(pkg)
	return payload, err
}

func (d *DefaultConnect) writeByte(bytes []byte) error {
	size := varint.New(int64(len(bytes)))
	_, err := d.c.Write(append(size, bytes...))
	return err
}

func (d *DefaultConnect) readByte() ([]byte, error) {
	if d.close {
		return nil, nil
	}
	reader := bufio.NewReader(d.c)
	i, err := binary.ReadVarint(reader)
	if err != nil {
		return nil, err
	}
	all, err := ioutil.ReadAll(io.LimitReader(reader, i))
	return all, err
}

func NewConnect(c net.Conn) *DefaultConnect {
	return &DefaultConnect{c: c}
}

func (d *DefaultConnect) Close() error {
	d.l.Lock()
	defer d.l.Unlock()
	if d.close {
		return nil
	}
	d.close = true
	return d.c.Close()
}
func (d *DefaultConnect) Closed() bool {
	d.l.Lock()
	defer d.l.Unlock()
	return d.close
}
