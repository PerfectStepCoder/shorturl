package config

import (
	"errors"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

func (a NetAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

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

type Settings struct {
	ServiceNetAddress NetAddress
	BaseURL           string
	FileStoragePath   string
	DatabaseDSN       string
	SaveDBtoFile      bool
}
