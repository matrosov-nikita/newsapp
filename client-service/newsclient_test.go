package client_service

import (
	"errors"
	"testing"

	"github.com/golang/protobuf/ptypes"

	pb "github.com/matrosov-nikita/newsapp/client-service/proto"
	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	suite.Suite

	client *NewsClient
	mq     *SpyMQSender
}

func (s *ClientSuite) SetupTest() {
	s.mq = &SpyMQSender{}
	s.client = New(s.mq)
}

func (s *ClientSuite) TestRequestOfMessageQueueFails() {
	s.mq.ReturnErrors(true)
	_, err := s.client.CreateNews("header")
	s.NotNil(err)
}

func (s *ClientSuite) TestMQRequestCalledWithCorrectData() {
	_, err := s.client.CreateNews("header")
	s.Nil(err)
	s.Equal("header", s.mq.bodyData.Header)
}

func (s *ClientSuite) TestErrorCodeGraterThanZeroWhenCreate() {
	s.mq.ConfigureResponse(1, "some error")
	_, err := s.client.CreateNews("header")
	s.Equal(errors.New("some error"), err)
}

func (s *ClientSuite) TestMQFindRequestFails() {
	s.mq.ReturnErrors(true)
	_, err := s.client.FindById("id")
	s.NotNil(err)
}

func (s *ClientSuite) TestErrorCodeGraterThanZeroWhenFind() {
	s.mq.ConfigureResponse(1, "some error")
	_, err := s.client.FindById("id")
	s.Equal(errors.New("some error"), err)
}

func (s *ClientSuite) TestNotFoundError() {
	s.mq.ConfigureResponse(2, "news not found")
	_, err := s.client.FindById("id")
	s.Equal(ErrNewsNotFound, err)
}

func (s *ClientSuite) TestFindNewsReturnResponse() {
	resp, err := s.client.FindById("id")
	s.Nil(err)
	s.Equal("123", resp.ID)
	s.Equal("header", resp.Header)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

type SpyMQSender struct {
	bodyData     *pb.News
	returnErrors bool

	errorCode int32
	error     string
}

func (r *SpyMQSender) ConfigureResponse(code int32, error string) {
	r.errorCode = code
	r.error = error
}

func (r *SpyMQSender) Find(data *pb.FindRequest) (*pb.FindResponse, error) {
	if r.returnErrors {
		return nil, errors.New("message queue error")
	}

	return &pb.FindResponse{
		News: &pb.News{
			Id:        "123",
			Header:    "header",
			CreatedAt: ptypes.TimestampNow(),
		},
		ErrorCode: r.errorCode,
		Error:     r.error,
	}, nil
}

func (r *SpyMQSender) ReturnErrors(v bool) {
	r.returnErrors = v
}

func (r *SpyMQSender) Create(data *pb.News) (*pb.CreateResponse, error) {
	if r.returnErrors {
		return nil, errors.New("message queue error")
	}

	r.bodyData = data
	return &pb.CreateResponse{
		ErrorCode: r.errorCode,
		Error:     r.error,
	}, nil
}
