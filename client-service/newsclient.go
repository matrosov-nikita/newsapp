package client_service

import (
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/matrosov-nikita/newsapp/client-service/proto"
)

// ErrNewsNotFound happens when news can not be find by id.
var ErrNewsNotFound = errors.New("news not found")

type MQSender interface {
	Create(data *pb.News) (*pb.CreateResponse, error)
	Find(data *pb.FindRequest) (*pb.FindResponse, error)
}

// NewsClient represents client service.
type NewsClient struct {
	mq MQSender
}

func New(m MQSender) *NewsClient {
	return &NewsClient{mq: m}
}

// CreateNews creates news by given header.
func (c *NewsClient) CreateNews(header string) (string, error) {
	resp, err := c.mq.Create(&pb.News{
		Header:    header,
		CreatedAt: ptypes.TimestampNow(),
	})

	if err != nil {
		return "", err
	}

	if resp.ErrorCode > 0 {
		return "", errors.New(resp.Error)
	}

	return resp.Id, nil
}

type ResponseNews struct {
	ID        string    `json:"id"`
	Header    string    `json:"header"`
	CreatedAt time.Time `json:"createdAt"`
}

// FindById finds news by given id.
func (c *NewsClient) FindById(id string) (*ResponseNews, error) {
	resp, err := c.mq.Find(&pb.FindRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	if resp.ErrorCode > 0 {
		if resp.ErrorCode == 2 {
			return nil, ErrNewsNotFound
		}

		return nil, errors.New(resp.Error)
	}

	createdAt, _ := ptypes.Timestamp(resp.News.CreatedAt)
	return &ResponseNews{
		ID:        resp.News.Id,
		Header:    resp.News.Header,
		CreatedAt: createdAt,
	}, nil
}
