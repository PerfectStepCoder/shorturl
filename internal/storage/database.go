package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type StorageInPostgres struct {
	connectionToDB *pgx.Conn
	lengthShortURL int
}

func initDB(config *pgx.ConnConfig) bool {
	// Подключение к стандартной БД
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		"postgres")
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return false
	}
	defer conn.Close(context.Background())
	// Создание базы данных, если ее нет
	result, err := conn.Exec(context.Background(), "CREATE DATABASE urlservice")
	log.Print(result)
	if err != nil {
		// Проверка на ошибку, если база данных уже существует
		log.Printf("Database already exist: %v\n", err)
	}
	// Подключение к новой базе данных "urlservice"
	connString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database)
	urlserviceDB, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return false
	}
	defer urlserviceDB.Close(context.Background())

	// Создание таблицы "urls", если ее нет
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		uuid UUID PRIMARY KEY,
		short VARCHAR(255) NOT NULL,
		original TEXT NOT NULL
	)`
	_, err = urlserviceDB.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Failed to create table: %v\n", err)
		return false
	}
	return true
}

func NewStorageInPostgres(connectionString string, lengthShortURL int) (*StorageInPostgres, error) {

	newStorage := StorageInPostgres{connectionToDB: nil, lengthShortURL: lengthShortURL}

	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		log.Printf("Failed to parse connection string: %v", err)
		return &newStorage, errors.New("Failed to connect to database")
	}

	if initDB(config) {
		connectionToDB, err := pgx.Connect(context.Background(), connectionString)
		if err != nil {
			return &newStorage, errors.New("Failed to connect to database")
		}
		newStorage.connectionToDB = connectionToDB
		return &newStorage, nil
	} else {
		return &newStorage, errors.New("Failed database init")
	}
}

func (s *StorageInPostgres) Get(hashKey string) (string, bool) {
	var originalURL string
	query := `
		SELECT original FROM urls WHERE short = $1
	`
	err := s.connectionToDB.QueryRow(context.Background(), query, hashKey).Scan(&originalURL)
	if err != nil {
		log.Printf("Failed to find original URL: %v\n", err)
		return originalURL, false
	}
	return originalURL, true
}

func (s *StorageInPostgres) Save(value string) string {
	newUUID := uuid.New()
	hashKey := make_hash(value, s.lengthShortURL)
	// SQL-запрос на вставку новой записи
	query := `
		INSERT INTO urls (uuid, short, original)
		VALUES ($1, $2, $3)
	`
	_, err := s.connectionToDB.Exec(context.Background(), query, newUUID, hashKey, value)

	if err != nil {
		log.Printf("Failed to insert new record: %v\n", err)
	}

	return hashKey
}

func (s *StorageInPostgres) LoadData(pathToFile string) int {
	// TODO реализовать загрузку БД из дампа
	return 0
}

func (s *StorageInPostgres) SaveData(pathToFile string) int {
	// TODO реализовать сохранение БД в дамп
	return 0
}

func (s *StorageInPostgres) Close() {
	s.connectionToDB.Close(context.Background())
}
