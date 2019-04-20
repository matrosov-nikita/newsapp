package storage_service

import (
	"errors"

	pb "github.com/matrosov-nikita/newsapp/storage-service/proto"
)

// ErrNewsNotFound happens when new cannot be found in storage.
var ErrNewsNotFound = errors.New("news not found")

type NewsStorage interface {
	Save(news *pb.News) (string, error)
	FindByID(id string) (*pb.News, error)
}

// StorageService represents service for storing news.
type StorageService struct {
	storage NewsStorage
}

// NewStorageService creates a new storage service.
func NewStorageService(storage NewsStorage) *StorageService {
	return &StorageService{storage: storage}
}

// Create saves news in storage.
func (s *StorageService) Create(news *pb.News) (*pb.CreateResponse, error) {
	id, err := s.storage.Save(news)

	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{
		Id: id,
	}, nil
}

// FindById finds new in storage by given id.
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
