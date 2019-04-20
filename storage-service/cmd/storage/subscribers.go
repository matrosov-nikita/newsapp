package main

import (
	"fmt"
	"log"

	"github.com/matrosov-nikita/newsapp/storage-service"

	"github.com/golang/protobuf/proto"
	pb "github.com/matrosov-nikita/newsapp/storage-service/proto"
	"github.com/nats-io/go-nats"
)

type Subs struct {
	st   *storage_service.StorageService
	conn *nats.Conn
}

func NewSubs(st *storage_service.StorageService, conn *nats.Conn) *Subs {
	return &Subs{
		conn: conn,
		st:   st,
	}
}

func (s *Subs) CreateNews(m *nats.Msg) {
	var news pb.News

	err := proto.Unmarshal(m.Data, &news)
	if err != nil {
		fail(s.conn, m, 1, fmt.Sprintf("cannot unmarshal request: %v", err))
		return
	}

	resp, err := s.st.Create(&news)
	if err != nil {
		fail(s.conn, m, 1, fmt.Sprintf("cannot save news: %v", err))
		return
	}

	bs, err := proto.Marshal(resp)
	if err == nil {
		if err := s.conn.Publish(m.Reply, bs); err != nil {
			log.Printf("could not publish message to queue: %v\n", err)
		}
	}
}

func (s *Subs) FindNews(m *nats.Msg) {
	var reqID pb.FindRequest

	err := proto.Unmarshal(m.Data, &reqID)
	if err != nil {
		fail(s.conn, m, 1, fmt.Sprintf("cannot unmarshal request: %v", err))
		return
	}

	resp, err := s.st.FindById(&reqID)
	if err != nil {
		fail(s.conn, m, 1, fmt.Sprintf("cannot find news: %v", err))
		return
	}

	bs, err := proto.Marshal(resp)
	if err == nil {
		if err := s.conn.Publish(m.Reply, bs); err != nil {
			log.Printf("could not publish message to queue: %v\n", err)
		}
	}
}

func fail(conn *nats.Conn, m *nats.Msg, errorCode int32, error string) {
	resp := pb.CreateResponse{
		ErrorCode: errorCode,
		Error:     error,
	}

	bs, err := proto.Marshal(&resp)
	if err == nil {
		conn.Publish(m.Reply, bs)
	}
}
