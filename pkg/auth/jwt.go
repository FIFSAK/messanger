package auth

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
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
func CreateToken(username string, id int, writer http.ResponseWriter) error {
	payload := jwt.MapClaims{
		"sub": username,
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	payloadRefresh := jwt.MapClaims{
		"sub": username,
		"id":  id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, payloadRefresh)
	t, err := token.SignedString([]byte(os.Getenv("secretKey")))
	if err != nil {
		http.Error(writer, "jwt token signing", http.StatusBadRequest)
	}
	tr, err := tokenRefresh.SignedString([]byte(os.Getenv("secretKey")))
	if err != nil {
		http.Error(writer, "jwt refresh token signing", http.StatusBadRequest)
	}

	err = json.NewEncoder(writer).Encode(struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
	}{
		Token:        t,
		RefreshToken: tr,
	})
	if err != nil {
		return err
	}
	return nil
}
