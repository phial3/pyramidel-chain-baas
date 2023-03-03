package mq

import (
	"github.com/streadway/amqp"
	"os"
	"reflect"
	"testing"
)

func TestNewRabbitMQ(t *testing.T) {
	tests := []struct {
		name string
		want *RabbitMQ
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRabbitMQ(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRabbitMQ() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRabbitMQ_Connect(t *testing.T) {
	type fields struct {
		conn      *amqp.Connection
		channel   *amqp.Channel
		queue     *amqp.Queue
		QueueName string
		Exchange  string
		Key       string
		Mqurl     string
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
			mq := &RabbitMQ{
				conn:      tt.fields.conn,
				channel:   tt.fields.channel,
				queue:     tt.fields.queue,
				QueueName: tt.fields.QueueName,
				Exchange:  tt.fields.Exchange,
				Key:       tt.fields.Key,
				Mqurl:     tt.fields.Mqurl,
			}
			if err := mq.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRabbitMQ_Destroy(t *testing.T) {
	type fields struct {
		conn      *amqp.Connection
		channel   *amqp.Channel
		queue     *amqp.Queue
		QueueName string
		Exchange  string
		Key       string
		Mqurl     string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mq := &RabbitMQ{
				conn:      tt.fields.conn,
				channel:   tt.fields.channel,
				queue:     tt.fields.queue,
				QueueName: tt.fields.QueueName,
				Exchange:  tt.fields.Exchange,
				Key:       tt.fields.Key,
				Mqurl:     tt.fields.Mqurl,
			}
			mq.Destroy()
		})
	}
}

func TestRabbitMQ_Publish(t *testing.T) {
	if err := os.Setenv("PYCBAAS_CFG_PATH", "E:\\github.com\\hxx258456\\pyramidel-chain-baas\\configs"); err != nil {
		panic(err)
	}
	type fields struct {
		conn      *amqp.Connection
		channel   *amqp.Channel
		queue     *amqp.Queue
		QueueName string
		Exchange  string
		Key       string
		Mqurl     string
	}
	type args struct {
		message []byte
	}
	message := `
[
    {
        "Domain": "orderer9.91610136MAB0QXTC21.pcb.com",
        "dueTime": "2006-01-02 15:04:05",
        "restartTime": "2006-01-02 15:04:05",
        "nodeCore": 1,
        "nodeMemory": 1,
        "nodeBandwidth": 1,
        "nodeDisk": 10,
        "hostId": 10,
        "serialNumber": 9,
        "port": 23823,
        "name": "91610136MAB0QXTC21_orderer9",
        "organizationId": 1,
        "orgPackageId": 9,
        "status": 1,
        "error": "",
        "ID": 9,
        "CreatedAt": "2006-01-02 15:04:05",
        "UpdatedAt": "2006-01-02 15:04:05",
        "DeletedAt": [
            {
                "Domain": "orderer9.91610136MAB0QXTC21.pcb.com",
                "dueTime": "2006-01-02 15:04:05",
                "restartTime": "2006-01-02 15:04:05",
                "nodeCore": 1,
                "nodeMemory": 1,
                "nodeBandwidth": 1,
                "nodeDisk": 10,
                "hostId": 10,
                "serialNumber": 9,
                "port": 23823,
                "name": "91610136MAB0QXTC21_orderer9",
                "organizationId": 1,
                "orgPackageId": 9,
                "status": 1,
                "error": "",
                "ID": 9,
                "CreatedAt": "2006-01-02 15:04:05",
                "UpdatedAt": "2006-01-02 15:04:05",
                "DeletedAt": "2006-01-02 15:04:05"
            },
            {
                "Domain": "peer18.91610136MAB0QXTC21.pcb.com",
                "dueTime": "2006-01-02 15:04:05",
                "restartTime": "2006-01-02 15:04:05",
                "nodeCore": 1,
                "nodeMemory": 1,
                "nodeBandwidth": 2,
                "nodeDisk": 20,
                "hostId": 6,
                "port": 28899,
                "name": "91610136MAB0QXTC21_peer18",
                "serialNumber": 18,
                "organizationId": 1,
                "orgPackageId": 9,
                "status": 1,
                "ccPort": 23739,
                "dbPort": 26297,
                "error": "",
                "ID": 18,
                "CreatedAt": "2006-01-02 15:04:05",
                "UpdatedAt": "2006-01-02 15:04:05",
                "DeletedAt": "2006-01-02 15:04:05"
            },
            {
                "Domain": "peer17.91610136MAB0QXTC21.pcb.com",
                "dueTime": "2006-01-02 15:04:05",
                "restartTime": "2006-01-02 15:04:05",
                "nodeCore": 1,
                "nodeMemory": 1,
                "nodeBandwidth": 2,
                "nodeDisk": 20,
                "hostId": 2,
                "port": 30321,
                "name": "91610136MAB0QXTC21_peer17",
                "serialNumber": 17,
                "organizationId": 1,
                "orgPackageId": 9,
                "status": 1,
                "ccPort": 2377,
                "dbPort": 3109,
                "error": "",
                "ID": 17,
                "CreatedAt": "2006-01-02 15:04:05",
                "UpdatedAt": "2006-01-02 15:04:05",
                "DeletedAt": "2006-01-02 15:04:05"
            }
        ]
    },
    {
        "Domain": "peer18.91610136MAB0QXTC21.pcb.com",
        "dueTime": "2006-01-02 15:04:05",
        "restartTime": "2006-01-02 15:04:05",
        "nodeCore": 1,
        "nodeMemory": 1,
        "nodeBandwidth": 2,
        "nodeDisk": 20,
        "hostId": 6,
        "port": 28899,
        "name": "91610136MAB0QXTC21_peer18",
        "serialNumber": 18,
        "organizationId": 1,
        "orgPackageId": 9,
        "status": 1,
        "ccPort": 23739,
        "dbPort": 26297,
        "error": "",
        "ID": 18,
        "CreatedAt": "2006-01-02 15:04:05",
        "UpdatedAt": "2006-01-02 15:04:05",
        "DeletedAt": null
    },
    {
        "Domain": "peer17.91610136MAB0QXTC21.pcb.com",
        "dueTime": "2006-01-02 15:04:05",
        "restartTime": "2006-01-02 15:04:05",
        "nodeCore": 1,
        "nodeMemory": 1,
        "nodeBandwidth": 2,
        "nodeDisk": 20,
        "hostId": 2,
        "port": 30321,
        "name": "91610136MAB0QXTC21_peer17",
        "serialNumber": 17,
        "organizationId": 1,
        "orgPackageId": 9,
        "status": 1,
        "ccPort": 2377,
        "dbPort": 3109,
        "error": "",
        "ID": 17,
        "CreatedAt": "2006-01-02 15:04:05",
        "UpdatedAt": "2006-01-02 15:04:05",
        "DeletedAt": "2006-01-02 15:04:05"
    }
]
`
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				conn:      nil,
				channel:   nil,
				queue:     nil,
				QueueName: "baasOrgAdd",
				Exchange:  "",
				Key:       "",
				Mqurl:     "amqp://txhy:txhy2022.com@47.92.54.239:5672//auth",
			},
			args: args{
				message: []byte(message),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RabbitMQ{
				conn:      tt.fields.conn,
				channel:   tt.fields.channel,
				queue:     tt.fields.queue,
				QueueName: tt.fields.QueueName,
				Exchange:  tt.fields.Exchange,
				Key:       tt.fields.Key,
				Mqurl:     tt.fields.Mqurl,
			}
			if err := r.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := r.Publish(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
