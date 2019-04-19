package nats

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	pb "github.com/matrosov-nikita/newsapp/newsclient/proto"
	"github.com/nats-io/go-nats"
)

const (
	CreateSubject = "news.create"
	GetSubject    = "news.get"

	TimeoutMessageReceive = time.Second
)

type NatsClient struct {
	conn *nats.Conn
}

func New(natsURL string) (*NatsClient, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	return &NatsClient{
		conn: conn,
	}, nil
}

func (c *NatsClient) Close() {
	c.conn.Close()
}

func (c *NatsClient) Create(data *pb.News) (*pb.CreateResponse, error) {
	bs, err := proto.Marshal(data)
	if err != nil {
		return nil, err
	}

	msg, err := c.conn.Request(CreateSubject, bs, TimeoutMessageReceive)
	if err != nil || msg == nil {
		return nil, fmt.Errorf("request error or message is empty: %v", err)
	}
	var resp pb.CreateResponse
	err = proto.Unmarshal(msg.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error: %v", err)
	}

	return &resp, nil
}

func (c *NatsClient) Find(id string) (*pb.FindResponse, error) {
	bs, err := proto.Marshal(&pb.FindRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	msg, err := c.conn.Request(GetSubject, bs, TimeoutMessageReceive)
	if err != nil || msg == nil {
		return nil, fmt.Errorf("request error or message is empty: %v", err)
	}
	var resp pb.FindResponse
	err = proto.Unmarshal(msg.Data, &resp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error: %v", err)
	}

	return &resp, nil
}
