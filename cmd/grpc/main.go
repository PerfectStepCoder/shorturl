package main

import (
	"fmt"
	"log"
	"net"

	"github.com/PerfectStepCoder/shorturl/cmd/shortener/config"
	pb "github.com/PerfectStepCoder/shorturl/internal/proto/gen"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"google.golang.org/grpc"
)

// mainStorage - хранилище для записи и чтения обработанных ссылок.
var mainStorage storage.PersistanceStorage

const (
	// lengthShortURL — константа длина генерируемых коротких ссылок.
	lengthShortURL = 10
)

func main() {

	var logger, logFile = config.GetLogger()
	defer logFile.Close()

	appSettings, mainStorage, errConfig := config.SetupServerSettings(logger, lengthShortURL)

	fmt.Println(appSettings)

	if errConfig != nil {
		panic("No work without correct config!")
	}

	defer mainStorage.Close()
	
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", appSettings.ServiceNetAddress.Port))
	if err != nil {
		log.Fatal(err)
	}

	// Определяем методы, которые будут игнорироваться перехватчиком.
	ignoredMethods := []string{
		"Login", "Stats",
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptorAuth(ignoredMethods)))

	// регистрируем сервис
	pb.RegisterShorterServer(s, &ShorterServerGRPC{
		mainStorage: mainStorage,
	})

	log.Print("Сервер gRPC начал работу")

	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}