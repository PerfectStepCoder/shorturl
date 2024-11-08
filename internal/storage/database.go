// Модуль содержит реализацию интерфейса хранилища для БД Postgres
package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBPool - интерфейс для пула соеденений
type DBPool interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Close()
}

// StorageInMemory - хранилище в базе данных Postgres
type StorageInPostgres struct {
	connectionToDB     *pgx.Conn
	poolConnectionToDB DBPool // Используем пул соединений *pgxpool.Pool
	lengthShortURL     int
}

var cache map[string]bool = make(map[string]bool)

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
		original TEXT NOT NULL UNIQUE,
		user_uid VARCHAR(1024) NULL,
		deleted BOOLEAN DEFAULT false     
	)`
	// ; CREATE INDEX idx_user_uid ON urls (user_uid);
	// Тесты 11 проваливаются из за создани индекса
	// ;
	//-- Добавляем индекс на поле short
	//CREATE INDEX idx_short ON urls (short);

	_, err = urlserviceDB.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Failed to create table: %v\n", err)
		return false
	}
	return true
}

// NewStorageInPostgres - конструктор
func NewStorageInPostgres(connectionString string, lengthShortURL int) (*StorageInPostgres, error) {

	newStorage := StorageInPostgres{connectionToDB: nil, lengthShortURL: lengthShortURL}

	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		log.Printf("Failed to parse connection string: %v", err)
		return &newStorage, errors.New("failed to connect to database")
	}

	if initDB(config) {
		connectionToDB, err := pgx.Connect(context.Background(), connectionString)
		poolConfig, _ := pgxpool.ParseConfig(connectionString)
		poolConfig.MaxConns = 100
		poolConfig.MinConns = 50
		poolConnectionToDB, err1 := pgxpool.NewWithConfig(context.Background(), poolConfig)

		if err != nil || err1 != nil {
			return &newStorage, errors.New("failed to connect to database")
		}
		newStorage.connectionToDB = connectionToDB
		newStorage.poolConnectionToDB = poolConnectionToDB
		return &newStorage, nil
	} else {
		return &newStorage, errors.New("failed database init")
	}
}

// Get - чтение ссылки.
func (s *StorageInPostgres) Get(hashKey string) (string, bool) {
	var originalURL string

	query := "SELECT original FROM urls WHERE short = $1"

	err := s.poolConnectionToDB.QueryRow(context.Background(), query, hashKey).Scan(&originalURL)
	if err != nil {
		log.Printf("Failed to find original URL: %v\n", err)
		return originalURL, false
	}
	return originalURL, true
}

// IsDeleted - удалена ли ссылка.
func (s *StorageInPostgres) IsDeleted(hashKey string) (bool, error) {
	_, exists := cache[hashKey]
	return exists, nil
}

// Save - сохранение новой ссылки.
func (s *StorageInPostgres) Save(value string, userUID string) (string, error) {
	newUUID := uuid.New()
	hashKey := makeHash(value, s.lengthShortURL)
	// SQL-запрос на вставку новой записи
	query := `
		INSERT INTO urls (uuid, short, original, user_uid)
		VALUES ($1, $2, $3, $4)
	`
	_, err := s.poolConnectionToDB.Exec(context.Background(), query, newUUID, hashKey, value, userUID)

	if err != nil {
		// Проверка на ошибку типа UniqueViolation
		var pge *pgconn.PgError
		if errors.As(err, &pge) {
			if pge.Code == pgerrcode.UniqueViolation {
				log.Printf("Error: A url with the same value already exists. URL: %s, hash: %s", value, hashKey)
				return hashKey, NewUniqURLError(value, hashKey)
			}
		}
		log.Printf("Failed to insert new record: %v\n", err)
		return hashKey, err
	}
	return hashKey, nil
}

// FindByUserUID - поиск ссылок по пользовательскому UID.
func (s *StorageInPostgres) FindByUserUID(userUID string) ([]ShortHashURL, error) {
	var output []ShortHashURL
	// SQL-запрос на поиск URLs
	query := `
		SELECT short, original FROM urls WHERE user_uid = $1
	`
	urls, err := s.connectionToDB.Query(context.Background(), query, userUID)

	if err != nil {
		log.Printf("Failed to find original URL: %v\n", err)
		return output, err
	}

	defer urls.Close()

	// Итерируем по строкам результата
	for urls.Next() {
		var shortURL, originalURL string

		// Чтение данных в переменные
		err = urls.Scan(&shortURL, &originalURL)
		if err != nil {
			log.Printf("failed to scan row: %s", err)
			return output, err
		}

		// Добавление URL в массив
		output = append(output, ShortHashURL{
			ShortHash:   shortURL,
			OriginalURL: originalURL,
		})
	}

	if urls.Err() != nil {
		log.Printf("error after iterating rows: %s", urls.Err())
	}

	return output, nil
}

// DeleteByUser - удалить ссылку по пользовательскому UUID
func (s *StorageInPostgres) DeleteByUser(shortsHashURL []string, userUID string) error {

	// В кеш
	for _, v := range shortsHashURL { // записываем удаляемый батч в кеш, для запроса GET /{id}, чтобы не обращатся к БД
		cache[v] = true
	}

	// Создаем объект Batch
	batch := &pgx.Batch{}

	for _, shortHashURL := range shortsHashURL { // short - короткая ссылка
		batch.Queue("UPDATE urls SET deleted = true WHERE short = $1 and user_uid = $2", shortHashURL, userUID)
	}

	batchResults := s.poolConnectionToDB.SendBatch(context.Background(), batch)
	defer batchResults.Close() // Закрываем BatchResults после использования

	// Обработка каждой команды в батче
	for i := 0; i < len(shortsHashURL); i++ {
		_, err := batchResults.Exec()
		if err != nil {
			log.Printf("Error executing batch command: %v", err)
			return err
		}
	}

	return nil
}

// LoadData загрузка данных из файла
func (s *StorageInPostgres) LoadData(pathToFile string) int {
	// TODO реализовать загрузку БД из дампа
	return 0
}

// SaveData сохранение данных в файл
func (s *StorageInPostgres) SaveData(pathToFile string) int {
	// TODO реализовать сохранение БД в дамп
	return 0
}

// Close - освобождение ресурсов
func (s *StorageInPostgres) Close() {
	s.connectionToDB.Close(context.Background())
}

// CorrelationSave - сохранение данных (ссылка и идентификатор)
func (s *StorageInPostgres) CorrelationSave(value string, correlationID string, userUID string) string {
	// SQL-запрос на вставку новой записи
	query := `
		INSERT INTO urls (uuid, short, original, user_uid)
		VALUES ($1, $2, $3, $4)
	`
	_, err := s.connectionToDB.Exec(context.Background(), query, correlationID, correlationID, value, userUID)

	if err != nil {
		log.Printf("Failed to insert new record: %v\n", err)
	}

	return correlationID
}

// CorrelationGet - чтение данных (ссылка и идентификатор)
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

// CorrelationSave - сохранение данных (ссылок и идентификатор)
func (s *StorageInPostgres) CorrelationsSave(correlationURLs []CorrelationURL, userUID string) ([]string, error) {

	var output []string

	// Подготовка SQL-запроса для вставки данных
	query := `INSERT INTO urls (uuid, correlation_id, short, original, user_uid) VALUES ($1, $2, $3, $4, $5)`

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
		_, err = tx.Exec(context.Background(), query, newUUID, item.CorrelationID, shortURL, originalURL, userUID)
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
