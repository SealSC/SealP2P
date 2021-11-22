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
	"encoding/base64"
)

type Message struct {
	Family    string
	Version   string
	Action    string
	Payload   base64Byte
	Hash      base64Byte
	Signature base64Byte

	FromID string
	TS     int64
	ToID   []string
}
type base64Byte []byte

func (b base64Byte) String() string {
	return base64.StdEncoding.EncodeToString(b)
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

func (m *Message) UNPackByte(payload []byte) (err error) {
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
	s := strings.TrimSpace(string(split[4]))
	if len(s) != 0 {
		m.ToID = strings.Split(s, ";")
	}
	m.Hash, err = base64.StdEncoding.DecodeString(string(split[5]))
	if err != nil {
		return
	}
	if len(body) > 0 {
		m.Payload, err = base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			return
		}
	}
	return nil
}
