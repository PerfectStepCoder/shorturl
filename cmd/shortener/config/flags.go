package config

import (
    "flag"
)

func ParseFlags() Settings {
    appSettings := new(Settings)
	// Default
	appSettings.ServiceNetAddress.Host = "localhost"
	appSettings.ServiceNetAddress.Port = 8080
	flag.Var(&appSettings.ServiceNetAddress, "a", "Net address host:port")
	flag.StringVar(&appSettings.BaseUrl, "b", "http://localhost:8080", "Base url host:port") 
    flag.Parse()
    return *appSettings
}
