package container

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"reflect"
	"testing"
)

func TestNewContainerService(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name string
		args args
		want ContainerService
	}{
		{"test", args{host: "39.100.110.117"}, &containerService{
			Host: "39.100.110.117",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContainerService(tt.args.host); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContainerService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_containerService_Close(t *testing.T) {
	type fields struct {
		Host string
		cli  *client.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &containerService{
				Host: tt.fields.Host,
				cli:  tt.fields.cli,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_containerService_Conn(t *testing.T) {
	type fields struct {
		Host string
		cli  *client.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"success", fields{Host: "39.100.110.117", cli: nil}, true},
		{"error", fields{Host: "localhost", cli: nil}, true},
		{"internalip", fields{Host: "192.168.0.170", cli: nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &containerService{
				Host: tt.fields.Host,
				cli:  tt.fields.cli,
			}
			if err := s.Conn(); (err != nil && s.cli != nil) != tt.wantErr {
				t.Errorf("Conn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_containerService_List(t *testing.T) {
	type fields struct {
		Host string
		cli  *client.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []types.Container
		wantErr bool
	}{
		{"success", fields{Host: "39.100.110.117", cli: nil}, args{ctx: context.Background()}, []types.Container{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &containerService{
				Host: tt.fields.Host,
				cli:  tt.fields.cli,
			}
			err := s.Conn()
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := s.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_containerService_Login(t *testing.T) {
	type fields struct {
		Host string
		cli  *client.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &containerService{
				Host: tt.fields.Host,
				cli:  tt.fields.cli,
			}
			if err := s.Login(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_containerService_Pull(t *testing.T) {
	type fields struct {
		Host string
		cli  *client.Client
	}
	type args struct {
		ctx context.Context
		s2  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &containerService{
				Host: tt.fields.Host,
				cli:  tt.fields.cli,
			}
			if err := s.Pull(tt.args.ctx, tt.args.s2); (err != nil) != tt.wantErr {
				t.Errorf("Pull() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_containerService_Run(t *testing.T) {
	type fields struct {
		Host string
		cli  *client.Client
	}
	type args struct {
		ctx              context.Context
		config           *container.Config
		hostConfig       *container.HostConfig
		networkingConfig *network.NetworkingConfig
		containerName    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &containerService{
				Host: tt.fields.Host,
				cli:  tt.fields.cli,
			}
			if err := s.Run(tt.args.ctx, tt.args.config, tt.args.hostConfig, tt.args.networkingConfig, tt.args.containerName); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
