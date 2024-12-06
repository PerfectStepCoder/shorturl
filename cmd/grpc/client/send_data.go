package main

import (
	"context"
	"log"

	pb "github.com/PerfectStepCoder/shorturl/internal/proto/gen"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ParseError - парсинг ошибки приходящий от сервера.
func ParseError(err error) {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			log.Println(`NOT FOUND`, e.Message())
		default:
			// В остальных случаях выводим код ошибки в виде строки и сообщение
			log.Println(e.Code(), e.Message())
		}
	} else {
		log.Printf("Problem with parsing of error: %v", err)
	}
}

// TestURLs - передача тестовых данных к сервесу сокращения ссылок.
func TestURLs(c pb.ShorterClient) {

	// Набор тестовых данных
	urls := []*pb.RequestFullURL{
		{Url: faker.URL()},
		{Url: faker.URL()},
		{Url: faker.URL()},
		{Url: faker.URL()},
	}

	// X-Real-IP
	md := metadata.New(map[string]string{"X-Real-IP": "192.168.0.5"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Авторизация
	randomUID := uuid.New().String()
	token, err := c.Login(ctx, &pb.RequestLogin{
		UserUID: randomUID,
	})
	if err != nil {
		ParseError(err)
		return
	}

	// Авторизованные запросы к сервису
	md.Append("token", token.Jwt)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	log.Println(token)

	for _, url := range urls {
		resp, err := c.ShorterURL(ctx, url)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(resp.Result)
	}

	resultStats, _ := c.Stats(ctx, &emptypb.Empty{})

	log.Printf("\nStats:\nCountURL: %d\nCountUsers: %d", resultStats.CountURL, resultStats.CountUsers)

}
