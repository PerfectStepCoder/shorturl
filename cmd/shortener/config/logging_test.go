// Модуль config содержит настройки сервиса.
package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatLogFile(t *testing.T) {

	var _, logFile = GetLogger()
	defer logFile.Close()
	assert.FileExists(t, "logfile.log")

}
