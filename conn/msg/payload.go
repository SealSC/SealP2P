package msg

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"strconv"
	"crypto/sha1"
	"github.com/SealSC/SealP2P/tools/grand"
	"time"
)

type Payload struct {
	Version string
	FromID  string
	Path    string
	TS      int64
	ToID    []string

	MsgHash string
	Body    []byte
}

func (p *Payload) PackByte() []byte {
	if p.TS == 0 {
		p.TS = time.Now().Unix()
	}
	hash := sha1.New()
	int63 := grand.Int63()
	hash.Write([]byte(fmt.Sprintf("%s%d%d", p.FromID, p.TS, int63)))
	p.MsgHash = fmt.Sprintf("%x", hash.Sum(nil))
	tos := strings.Join(p.ToID, ";")
	head := fmt.Sprintf("%s|%s|%s|%d|%s|%s\n", p.Version, p.FromID, p.Path, p.TS, tos, p.MsgHash)
	payload := append([]byte(head), p.Body...)
	return payload
}

func (p *Payload) UNPackByte(payload []byte) error {
	index := bytes.IndexByte(payload, '\n')
	if index < 0 {
		return nil
	}
	head := payload[:index]
	body := payload[index+1:]
	split := bytes.Split(head, []byte("|"))
	if len(split) != 6 {
		return errors.New("payload head len != 6")
	}
	p.Version = string(split[0])
	p.FromID = string(split[1])
	p.Path = string(split[2])
	p.TS, _ = strconv.ParseInt(string(split[3]), 10, 64)
	p.ToID = strings.Split(string(split[4]), ";")
	p.MsgHash = string(split[5])
	p.Body = body
	return nil
}
