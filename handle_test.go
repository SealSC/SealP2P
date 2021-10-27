package SealP2P

import (
	"reflect"
	"testing"
	"github.com/SealSC/SealP2P/conn/msg"
)

func TestDefaultHandler_RegisterHandler(t *testing.T) {
	type args struct {
		key string
		f   func(req *msg.Message) *msg.Message
	}
	tests := []struct {
		name     string
		args     args
		watExist bool
	}{
		{name: "001", watExist: false, args: args{}},
		{name: "002", watExist: false, args: args{f: nil}},
		{name: "003", watExist: false, args: args{key: "", f: nil}},
		{name: "004", watExist: false, args: args{key: "kk", f: nil}},
		{name: "005", watExist: true, args: args{key: "kk", f: func(req *msg.Message) *msg.Message { return nil }}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultHandler{}
			d.RegisterHandler(tt.args.key, tt.args.f)
			f, ok := d.customMap[tt.args.key]
			exist := ok && f != nil
			if exist != tt.watExist {
				t.Errorf("RegisterHandler() = %v, want %v", exist, tt.watExist)
			}
		})
	}
}

type testMsg struct{}

func (t *testMsg) OnMessage(p *msg.Message) *msg.Message {
	return nil
}

func TestDefaultHandler_SetMessenger(t *testing.T) {

	type args struct {
		m Messenger
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "001", args: args{m: nil}},
		{name: "002", args: args{m: &testMsg{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultHandler{}
			d.SetMessenger(tt.args.m)
			if d.m != tt.args.m {
				t.Errorf("SetMessenger() m!=tt.args.m")
			}
		})
	}
}

func TestDefaultHandler_doHandle(t *testing.T) {
	type fields struct {
		customMap map[string]func(payload *msg.Message) *msg.Message
	}
	type args struct {
		req *msg.Message
	}
	xxxReq := NewPayload("xxx1")
	xxxReq.FromID = "xxx"
	xxxResp := NewPayload("xxx1")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *msg.Message
	}{
		{name: "001", fields: fields{customMap: nil}, args: args{req: nil}, want: nil},
		{name: "local_req", fields: fields{customMap: nil}, args: args{req: NewPayload(msg.Join)}, want: nil},
		{name: "003", fields: fields{customMap: map[string]func(payload *msg.Message) *msg.Message{"xxx1": func(*msg.Message) *msg.Message {
			return xxxResp
		}}}, args: args{req: xxxReq}, want: xxxResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultHandler{
				customMap: tt.fields.customMap,
			}
			if got := d.doHandle(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}
