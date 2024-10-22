// Модуль для создания логгера
package config

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// GetLogger - функция для создания логгера
func GetLogger() (*logrus.Logger, *os.File) {

	logger := logrus.New()

	// Открываем файл для логирования
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Error(err)
	}

	// Устанавливаем вывод логгера на несколько writer'ов (консоль и файл)
	multiWriter := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(multiWriter)

	// Устанавливаем формат вывода
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger, file
}
