package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func splitHostPort(addr string) (string, int, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid address format")
	}
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("can not str into number:", err)
		return parts[0], 0, nil
	}
	return parts[0], num, nil
}

func ParseFlags() Settings {
	appSettings := new(Settings)

	// Default
	appSettings.ServiceNetAddress.Host = ""
	appSettings.ServiceNetAddress.Port = 0
	baseURL := "http://localhost:8080"

	if envRunAddr := os.Getenv("SHORTURL_SERVER_ADDRESS"); envRunAddr != "" {
		host, port, err := splitHostPort(envRunAddr)
		if err == nil {
			appSettings.ServiceNetAddress.Host = host
			appSettings.ServiceNetAddress.Port = port
		}
	}

	if envBaseURL := os.Getenv("SHORTURL_BASE_URL"); envBaseURL != "" {
		baseURL = envBaseURL
	}

	flag.Var(&appSettings.ServiceNetAddress, "a", "Net address host:port")
	flag.StringVar(&appSettings.BaseURL, "b", baseURL, "Base url host:port")
	flag.Parse()

	if appSettings.ServiceNetAddress.Host == "" && appSettings.ServiceNetAddress.Port == 0 {
		appSettings.ServiceNetAddress.Host = "localhost"
		appSettings.ServiceNetAddress.Port = 8080
	}
	return *appSettings
}
