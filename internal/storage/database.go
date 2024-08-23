package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
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
		correlation_id TEXT NULL,
		short VARCHAR(255) NOT NULL,
		original TEXT NOT NULL UNIQUE
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
		return &newStorage, errors.New("failed to connect to database")
	}

	if initDB(config) {
		connectionToDB, err := pgx.Connect(context.Background(), connectionString)
		if err != nil {
			return &newStorage, errors.New("failed to connect to database")
		}
		newStorage.connectionToDB = connectionToDB
		return &newStorage, nil
	} else {
		return &newStorage, errors.New("failed database init")
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

func (s *StorageInPostgres) Save(value string) (string, error) {
	newUUID := uuid.New()
	hashKey := makeHash(value, s.lengthShortURL)
	// SQL-запрос на вставку новой записи
	query := `
		INSERT INTO urls (uuid, short, original)
		VALUES ($1, $2, $3)
	`
	_, err := s.connectionToDB.Exec(context.Background(), query, newUUID, hashKey, value)

	if err != nil {
		// Проверка на ошибку типа UniqueViolation
		var pge *pgconn.PgError
		if errors.As(err, &pge) {
			if pge.Code == pgerrcode.UniqueViolation {
				log.Println("Error: A url with the same value already exists.")
				return hashKey, NewUniqURLError(value, hashKey)
			}
		}
		log.Printf("Failed to insert new record: %v\n", err)
		return hashKey, err
	}
	return hashKey, nil
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

func (s *StorageInPostgres) CorrelationSave(value string, correlationID string) string {
	// SQL-запрос на вставку новой записи
	query := `
		INSERT INTO urls (uuid, short, original)
		VALUES ($1, $2, $3)
	`
	_, err := s.connectionToDB.Exec(context.Background(), query, correlationID, correlationID, value)

	if err != nil {
		log.Printf("Failed to insert new record: %v\n", err)
	}

	return correlationID
}

func (s *StorageInPostgres) CorrelationGet(correlationID string) (string, bool) {
	var originalURL string

	query := `
		SELECT original FROM urls WHERE short = $1
	`
	err := s.connectionToDB.QueryRow(context.Background(), query, correlationID).Scan(&originalURL)
	if err != nil {
		log.Printf("Failed to find original URL: %v\n", err)
		return originalURL, false
	}
	return originalURL, true
}

func (s *StorageInPostgres) CorrelationsSave(correlationURLs []CorrelationURL) ([]string, error) {

	var output []string

	// Подготовка SQL-запроса для вставки данных
	query := `INSERT INTO urls (uuid, correlation_id, short, original) VALUES ($1, $2, $3, $4)`

	// Начало транзакции
	tx, err := s.connectionToDB.Begin(context.Background())
	if err != nil {
		log.Printf("Failed to begin transaction: %v\n", err)
		return output, err
	}

	for _, item := range correlationURLs {

		newUUID := uuid.New()

		// Генерация UUID (например, с использованием pgx или других методов)
		shortURL := item.CorrelationID
		originalURL := item.OriginalURL
		output = append(output, shortURL)

		// Выполнение вставки в рамках транзакции
		_, err = tx.Exec(context.Background(), query, newUUID, item.CorrelationID, shortURL, originalURL)
		if err != nil {
			tx.Rollback(context.Background())
			log.Printf("Failed to insert data: %v\n", err)
			// Проверка на ошибку типа UniqueViolation
			var pge *pgconn.PgError
			if errors.As(err, &pge) {
				if pge.Code == pgerrcode.UniqueViolation {
					log.Println("Error: A url with the same value already exists.")
					return output, NewUniqURLError(originalURL, shortURL)
				}
			}
			return output, nil
		}
	}

	// Зафиксировать транзакцию
	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
	}

	return output, nil
}
