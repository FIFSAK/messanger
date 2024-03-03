package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

func JwtPayloadFromRequest(w http.ResponseWriter, r *http.Request) (jwt.MapClaims, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return nil, false
	}

	// Проверяем формат токена
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return nil, false
	}

	tokenString := authHeader[len(bearerPrefix):]
	jwtSecretKey := os.Getenv("secretKey")
	// Парсинг и валидация токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	fmt.Println(err)
	fmt.Println(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Токен валиден, и claims успешно извлечены
		return claims, true
	} else {
		// Токен невалиден или claims не могут быть приведены к типу MapClaims
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return nil, false
	}
}
