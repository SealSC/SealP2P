package msg

import (
	"testing"
	"reflect"
	"fmt"
	"time"
)

func TestMessage_Bae64(t *testing.T) {
	tests := []struct {
		name    string
		base64B base64Byte
		want    string
	}{
		{name: "empty", base64B: []byte(""), want: ""},
		{name: "132", base64B: []byte("123"), want: "MTIz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.base64B.String()
			if s != tt.want {
				t.Errorf("%s,!base64Byte", tt.name)
			}
		})
	}
}

func TestMessage_PackByte_UNPackByte(t *testing.T) {
	tests := []struct {
		name string
		msg  Message
	}{
		{name: "empty", msg: Message{}},
		{name: "time!=0", msg: Message{TS: time.Now().Unix()}},
		{name: "toID=1", msg: Message{ToID: []string{"0x11"}}},
		{name: "toID>1", msg: Message{ToID: []string{"0x12", "0x113"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := tt.msg.PackByte()
			temp := Message{}
			if err := temp.UNPackByte(bytes); err != nil {
				fmt.Println(string(bytes))
				fmt.Println()
				t.Errorf("%s,UNPackByte err:%v", tt.name, err)
			}
			if !reflect.DeepEqual(tt.msg, temp) {
				t.Errorf("%s,!reflect.DeepEqual", tt.name)
			}
		})
	}
}

func TestMessage_UNPackByte_Err(t *testing.T) {
	tests := []struct {
		name    string
		b       string
		wantErr bool
	}{
		{name: "is ok", b: "|||1637220857||PADbYpQz8n7ZFORq/FyWOuYme7o=\n", wantErr: false},
		{name: "body base64 decode err", b: "|||1637220857||PADbYpQz8n7ZFORq/FyWOuYme7o=\nxxx", wantErr: true},
		{name: "hash bae64 decode err", b: "|||1637220857||xixixixixi\n", wantErr: true},
		{name: "length<6", b: "|||163722085\n", wantErr: true},
		{name: "length>6", b: "|||||||||||||\n", wantErr: true},
		{name: "empty str", b: "", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := Message{}
			err := temp.UNPackByte([]byte(tt.b))
			if (err != nil) != tt.wantErr {
				t.Errorf("UNPackByte() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
