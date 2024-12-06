package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/PerfectStepCoder/shorturl/cmd/shortener/config"
	pb "github.com/PerfectStepCoder/shorturl/internal/proto/gen"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
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

	appSettings, mainStorage, errConfig := config.GetSettingsAndStorage(logger, lengthShortURL)

	log.Println(appSettings)

	if errConfig != nil {
		logger.Fatal("No work without correct config!")
	}

	defer mainStorage.Close()
	
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", appSettings.ServiceNetAddress.Port))
	if err != nil {
		logger.Fatal(err)
	}

	// Определяем методы, которые будут игнорироваться перехватчиком.
	ignoredMethods := map[string]struct{}{
		"/shorter.Shorter/Login": {}, "/shorter.Shorter/Stats": {},
	}

	// Создаём цепочку перехватчиков (порядок следования важен)
	unaryInterceptors := grpc_middleware.ChainUnaryServer(
		UnaryInterceptorTrustedSubnet(ignoredMethods, appSettings.TrustedSubnet),
		UnaryInterceptorAuth(ignoredMethods),                        
	)

	// Cоздаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptors))

	// Регистрируем сервис
	pb.RegisterShorterServer(s, &ShorterServerGRPC{
		mainStorage: mainStorage,
	})

	// Канал для обработки сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Print("Server gRPC start work")
		// Получаем запрос gRPC
		if err := s.Serve(listen); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	<-stop
	log.Print("Shutting down gRPC server...")

	// Корректное завершение работы
	gracefulStop(s)
	log.Print("Server stopped gracefully")

}

// Функция для graceful shutdown
func gracefulStop(server *grpc.Server) {
	
	// Завершаем сервер с ожиданием завершения текущих запросов
	server.GracefulStop()

}