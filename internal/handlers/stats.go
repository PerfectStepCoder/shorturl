// Модуль содержит декораторы для обработки запросов авторизованных пользователей.
package handlers

import (
	"encoding/json"
	"errors"
	"github.com/PerfectStepCoder/shorturl/internal/models"
	"github.com/PerfectStepCoder/shorturl/internal/storage"
	"log"
	"net"
	"net/http"
)

// ShorterStats - обрабатывает запрос к /api/internal/stats
func ShorterStats(mainStorage storage.Storage, trustedSubnet string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверка наличия trusted_subnet
		if trustedSubnet == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Получение IP из заголовка X-Real-IP
		clientIP := r.Header.Get("X-Real-IP")
		if clientIP == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Проверка, входит ли IP в доверенную подсеть
		isAllowed, err := isIPInCIDR(clientIP, trustedSubnet)
		if err != nil || !isAllowed {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Формирование и отправка ответа
		w.Header().Set("Content-Type", "application/json")

		users, _ := mainStorage.CountUsers()
		urls, _ := mainStorage.CountURLs()
		resp := models.ResponseStatsBase{
			Urls:  urls,
			Users: users,
		}

		// Cериализуем ответ сервера
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Println("Error writing response:", err)
			return
		}

		w.Write(jsonResp)
		w.WriteHeader(http.StatusOK)
	}
}

// isIPInCIDR проверяет, входит ли IP в подсеть
func isIPInCIDR(ip, cidr string) (bool, error) {
	_, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}

	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false, errors.New("invalid IP address")
	}

	return subnet.Contains(clientIP), nil
}
