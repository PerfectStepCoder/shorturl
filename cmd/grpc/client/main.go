package main

import (
	"fmt"
	"log"

	pb "github.com/PerfectStepCoder/shorturl/internal/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// serverPort - порт который слушает сервер
const serverPort = 8080

func main() {

	// Устанавливаем соединение с сервером
	conn, err := grpc.NewClient(fmt.Sprintf(":%d", serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Получаем переменную интерфейсного типа ShorterClient, через которую будем отправлять сообщения
	c := pb.NewShorterClient(conn)

	// функция, в которой будем отправлять сообщения
	TestURLs(c)
	
}