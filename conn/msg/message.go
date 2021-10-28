package msg

import (
	"time"
	"crypto/sha1"
	"github.com/SealSC/SealP2P/tools/grand"
	"fmt"
	"strings"
	"bytes"
	"errors"
	"strconv"
)

type Message struct {
	Family    string
	Version   string
	Action    string
	Payload   []byte
	Hash      []byte
	Signature []byte

	FromID string
	TS     int64
	ToID   []string
}

func (m *Message) PackByte() []byte {
	if m.TS == 0 {
		m.TS = time.Now().Unix()
	}
	hash := sha1.New()
	int63 := grand.Int63()
	hash.Write([]byte(fmt.Sprintf("%s%d%d", m.FromID, m.TS, int63)))
	m.Hash = hash.Sum(nil)
	tos := strings.Join(m.ToID, ";")
	head := fmt.Sprintf("%s|%s|%s|%d|%s|%s\n", m.Version, m.FromID, m.Action, m.TS, tos, m.Hash)
	payload := append([]byte(head), m.Payload...)
	return payload
}

func (m *Message) UNPackByte(payload []byte) error {
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
	m.Version = string(split[0])
	m.FromID = string(split[1])
	m.Action = string(split[2])
	m.TS, _ = strconv.ParseInt(string(split[3]), 10, 64)
	m.ToID = strings.Split(string(split[4]), ";")
	m.Hash = split[5]
	m.Payload = body
	return nil
}
