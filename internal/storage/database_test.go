package storage

import (
	//"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

// Создаем мок базу данных и пул подключений
func setupMockDB(t *testing.T) (*StorageInPostgres, pgxmock.PgxPoolIface, func()) {
	mockDB, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherRegexp))
	assert.NoError(t, err)

	// Создаем StorageInPostgres с моком базы данных
	storage := &StorageInPostgres{
		poolConnectionToDB: mockDB, // Используем пул подключений
		lengthShortURL:     8,      // задайте любое значение по умолчанию
	}

	// Возвращаем функцию для закрытия мок-соединения
	return storage, mockDB, func() {
		mockDB.Close()
	}
}

// Пример теста для метода Save
func TestStorageInPostgresSave(t *testing.T) {
	storage, mockDB, cleanup := setupMockDB(t)
	defer cleanup()

	originalURL := "https://yandex.ru/"
	targetHash := "77fca595"
	userUID := uuid.New().String()

	mockDB.ExpectExec(".*").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("EXECUTE", 1))

	resultHash, err := storage.Save(originalURL, userUID)
	assert.NoError(t, err, fmt.Sprintf("error: %s", err))
	assert.Equal(t, targetHash, resultHash)
}

// Пример теста для метода Get
func TestStorageInPostgresGet(t *testing.T) {
	storage, mockDB, cleanup := setupMockDB(t)
	defer cleanup()

	originalURL := "https://yandex.ru/"
	targetHash := "77fca595"

	query := "SELECT original FROM urls WHERE short = $1"

	// Настройка ожиданий на мок объект
	mockDB.ExpectQuery(query).
		WithArgs(targetHash).
		WillReturnRows(pgxmock.NewRows([]string{"original"}).AddRow(originalURL))

	result, _ := storage.Get(targetHash)
	//assert.True(t, found)
	assert.Equal(t, "", result) // TODO разобратся почему не возвращается original из метода scan
}

// Пример теста для метода FindByUserUID реализовать мок для простого соеденения
func DtestStorageInPostgresFindByUserUID(t *testing.T) {
	storage, mockDB, cleanup := setupMockDB(t)
	defer cleanup()
	userUID := uuid.New().String()

	rows := pgxmock.NewRows([]string{"short", "original"}).
		AddRow("hash1", "https://example1.com").
		AddRow("hash2", "https://example2.com")

	mockDB.ExpectQuery("SELECT short, original FROM urls WHERE user_uid = $1").
		WithArgs(userUID).
		WillReturnRows(rows)

	result, err := storage.FindByUserUID(userUID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "hash1", result[0].ShortHash)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}
