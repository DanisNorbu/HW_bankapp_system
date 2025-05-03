package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Проверка JWT токенов
func ValidateJWTToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение токена из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Разделение "Bearer <token>"
		tokenString := strings.Split(authHeader, " ")[1]

		// Проверка токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Мы проверяем только подпись, ключ у нас будет фиксирован
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем информацию о пользователе в контекст
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}

		clientID := claims["clientID"].(string)
		r.Header.Set("clientID", clientID)

		// Далее передаем выполнение следующему обработчику
		next(w, r)
	}
}

// Получение ID клиента из контекста запроса
func GetClientID(r *http.Request) string {
	return r.Header.Get("clientID")
}
