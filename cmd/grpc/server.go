package main

import (
	"context"
	"fmt"
    "google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/metadata" 
	pb "github.com/PerfectStepCoder/shorturl/internal/proto/gen"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
    "google.golang.org/protobuf/types/known/emptypb"
)

// UsersServer поддерживает все необходимые методы сервера.
type ShorterServerGRPC struct {
	pb.UnimplementedShorterServer
    mainStorage storage.PersistanceStorage
}

// ShorterURL - обработчик ссылок.
func (s *ShorterServerGRPC) ShorterURL(ctx context.Context, in *pb.RequestFullURL) (*pb.ResponseShortURL, error) {
    var response pb.ResponseShortURL
    var userUID string
    var md metadata.MD 
    var ok bool
    
    if md, ok = metadata.FromIncomingContext(ctx); !ok {
        return &response, status.Errorf(codes.Internal, `We didnot catch userUID`)
    }

    values := md.Get("userUID")
    if len(values) > 0 {
        userUID = values[0]
    }

	shortURL, err := s.mainStorage.Save(in.Url, userUID)
    if err != nil {
        return &response, status.Errorf(codes.Internal, `We didnot save your url`)
    }

	response.Result = shortURL

    fmt.Printf("Repite: %s, shortURL: %s\n", in.Url, shortURL)

	return &response, nil
}

// Login - логирование JWT токен.
func (s *ShorterServerGRPC) Login(ctx context.Context, in *pb.RequestLogin) (*pb.ResponseJWT, error) {
    var response pb.ResponseJWT
    var userUID string
    
    token, err := GenerateToken(userUID)

	if err != nil {
		return &response, status.Errorf(codes.Internal, `Problem to gen JWT`)
	}

    response.Jwt = token

    return &response, nil
}

// Stats - получение статситики, количество ссылок и пользователей
func (s *ShorterServerGRPC) Stats(ctx context.Context, in *emptypb.Empty) (*pb.ResponseStats, error) {
    var response pb.ResponseStats
    countURL, _ := s.mainStorage.CountURLs()
    countUsers, _ := s.mainStorage.CountUsers()
    response.CountURL = int32(countURL)
    response.CountUsers = int32(countUsers)
    return &response, nil
}
