package url

import (
	"context"
	"errors"
	usv1 "github.com/quibex/url-storage-api/gen/go/url-storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"url-storage/internal/cache"
)

type Storage interface {
	SetURL(ctx context.Context, url string, alias string) (int64, error)
	GetURL(ctx context.Context, alias string) (string, error)
}

type serverAPI struct {
	usv1.UnimplementedUrlStorageServer
	storage Storage
}

func Register(gRPC *grpc.Server, storage Storage) {
	usv1.RegisterUrlStorageServer(gRPC, &serverAPI{
		storage: storage,
	})
}

func (s *serverAPI) SetUrl(ctx context.Context, req *usv1.SetUrlRequest) (*usv1.SetUrlResponse, error) {
	if req.Url == "" || req.Alias == "" {
		return nil, status.Error(codes.InvalidArgument, "url and alias must not be empty")
	}

	id, err := s.storage.SetURL(ctx, req.Url, req.Alias)
	if err != nil {
		if errors.Is(err, cache.ErrAlreadyExist) {
			return nil, status.Error(codes.AlreadyExists, "alias already exist")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &usv1.SetUrlResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) GetUrl(ctx context.Context, req *usv1.GetUrlRequest) (*usv1.GetUrlResponse, error) {
	if req.Alias == "" {
		return nil, status.Error(codes.InvalidArgument, "alias must not be empty")
	}

	url, err := s.storage.GetURL(ctx, req.Alias)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "url not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &usv1.GetUrlResponse{
		Url: url,
	}, nil
}
