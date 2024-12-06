package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Секретный ключ для подписи токена.
var secretKey = []byte("your_secret_key")

// GenerateToken принимает userUID и возвращает JWT токен.
func GenerateToken(userUID string) (string, error) {
	// Создаем новый токен с данными.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userUID": userUID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Срок действия токена: 24 часа.
	})
	// Подписываем токен.
	return token.SignedString(secretKey)
}

// ExtractUserUID извлекает userUID из переданного токена.
func ExtractUserUID(tokenString string) (string, error) {
	// Парсим и проверяем токен.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	// Если парсинг не удался, возвращаем ошибку.
	if err != nil {
		return "", err
	}

	// Извлекаем данные из токена, если он валиден.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userUID, ok := claims["userUID"].(string); ok {
			return userUID, nil
		}
		return "", fmt.Errorf("userUID not found in token")
	}

	return "", fmt.Errorf("invalid token")
}
