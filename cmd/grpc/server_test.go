package main

import (
	"context"
	"testing"

	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	//"google.golang.org/grpc/status"
	pb "github.com/PerfectStepCoder/shorturl/internal/proto/gen"
)

const testLengthShortURL = 10

func TestShorterURL(t *testing.T) {

	inMemoryStorage, _ := storage.NewStorageInMemory(testLengthShortURL)
	defer inMemoryStorage.Close()

	targetURL := "http://example.com"
	
	// Создаём экземпляр сервиса
	server := &ShorterServerGRPC{
		mainStorage: inMemoryStorage,
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("userUID", "user_uuid_some"))

	// Создаём входной запрос
	request := &pb.RequestFullURL{Url: targetURL}

	// Вызываем метод для сокращения ссылок
	response, err := server.ShorterURL(ctx, request)

	// Проверяем результаты
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	savedURL, _ := server.mainStorage.Get(response.Result)
	if targetURL != savedURL {
		t.Errorf("expected %s, got %s", targetURL, response.Result)
	}

	// Вызываем метод статистики
	resultStats, _ := server.Stats(ctx, &emptypb.Empty{})
	if resultStats.CountURL != 1 {
		t.Errorf("expected CountURL = %d, got %d", 1, resultStats.CountURL)
	}
	if resultStats.CountUsers != 1 {
		t.Errorf("expected CountUsers = %d, got %d", 1, resultStats.CountUsers)
	}
}
