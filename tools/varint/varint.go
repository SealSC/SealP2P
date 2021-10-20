package varint

import (
	"encoding/binary"
)

func New(i int64) []byte {
	temp := make([]byte, 10)
	uvarint := binary.PutVarint(temp, i)
	return temp[:uvarint]
}
