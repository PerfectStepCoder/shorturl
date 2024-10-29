// Модуль config содержит настройки сервиса.
package config

import (
	"errors"
	"strconv"
	"strings"
)

// NetAddress - хост на котором будет доступен сервис.
type NetAddress struct {
	Host string
	Port int
}

// String - функция печати хоста сервиса.
func (a NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

// Set - функция установки хоста из передоваемой строки
func (a *NetAddress) Set(s string) error {
	if a.Host == "" && a.Port == 0 {
		hp := strings.Split(s, ":")
		if len(hp) != 2 {
			return errors.New("need address in a form host:port")
		}
		port, err := strconv.Atoi(hp[1])
		if err != nil {
			return err
		}
		a.Host = hp[0]
		a.Port = port
	}
	return nil
}

// Settings - все настройки сервиса.
type Settings struct {
	ServiceNetAddress NetAddress
	BaseURL           string
	FileStoragePath   string
	DatabaseDSN       string
	SaveDBtoFile      bool
	AddProfileRoute   bool
}
