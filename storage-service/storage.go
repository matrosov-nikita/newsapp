package storage_service

import (
	"errors"

	pb "github.com/matrosov-nikita/newsapp/storage-service/proto"
)

var ErrNewsNotFound = errors.New("news not found")

type NewsStorage interface {
	Save(news *pb.News) (string, error)
	FindByID(id string) (*pb.News, error)
}

type StorageService struct {
	storage NewsStorage
}

func NewStorageService(storage NewsStorage) *StorageService {
	return &StorageService{storage: storage}
}

func (s *StorageService) Create(news *pb.News) (*pb.CreateResponse, error) {
	id, err := s.storage.Save(news)

	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{
		Id: id,
	}, nil
}

func (s *StorageService) FindById(id string) (*pb.FindResponse, error) {
	news, err := s.storage.FindByID(id)
	if err != nil {
		if err == ErrNewsNotFound {
			return &pb.FindResponse{
				Error:     err.Error(),
				ErrorCode: 2,
			}, nil
		}

		return nil, err
	}

	return &pb.FindResponse{News: news}, nil
}
