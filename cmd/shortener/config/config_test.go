// Модуль config содержит настройки сервиса.
package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetAddress(t *testing.T) {

	netAddress := NetAddress{Host: "localhost", Port: 8080}
	assert.Equal(t, "localhost:8080", netAddress.String())

	// Set
	netAnotherAddress := NetAddress{Host: "", Port: 0}
	netAnotherAddress.Set("google:9999")
	assert.Equal(t, "google", netAnotherAddress.Host)
	assert.Equal(t, 9999, netAnotherAddress.Port)

}
