package SealP2P

import (
	"testing"
	"github.com/SealSC/SealP2P/conf"
)

func TestNewTcpService(t *testing.T) {
	type args struct {
		conf *conf.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "nil config", args: args{conf: nil}, wantErr: true},
		{name: "empty config", args: args{conf: &conf.Config{}}, wantErr: true},
		{name: "default config", args: args{conf: conf.DefaultConfig}, wantErr: true},
		{name: "default config exist node id", args: args{conf: &conf.Config{
			ID:            "xxxx",
			ClientOnly:    conf.DefaultConfig.ClientOnly,
			PKFile:        conf.DefaultConfig.PKFile,
			ServerPort:    conf.DefaultConfig.ServerPort,
			MulticastAddr: conf.DefaultConfig.MulticastAddr,
			MulticastPort: conf.DefaultConfig.MulticastPort,
		}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTcpService(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTcpService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTcpService_CloseAndDel(t1 *testing.T) {
	type fields struct {
		cache map[string]*ConnedNode
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "empty", fields: fields{map[string]*ConnedNode{}}, args: args{}},
		{name: "nil ", fields: fields{nil}, args: args{}},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TcpService{cache: tt.fields.cache}
			t.CloseAndDel(tt.args.key)
		})
	}
}
