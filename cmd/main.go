package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const telegramBotToken = "8248908460:AAFC0OhLPeu2r09gh_A7fecwMWErMCD5WBw"

// JWT секрет
const jwtSecret = "supersecretkey123"

// Разрешённый источник — именно HTTPS ngrok-домен фронта
const allowedOrigin = "https://tuyet-unaidable-elida.ngrok-free.dev"

// middleware CORS
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Проверка Telegram авторизации
func checkResponse(data map[string]string) bool {
	hashValue, ok := data["hash"]
	if !ok {
		return false
	}
	delete(data, "hash")

	keys := make([]string, 0, len(data))
	for k := range data {
		if data[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, data[k]))
	}
	dataString := strings.Join(pairs, "\n")

	secretKey := sha256.Sum256([]byte(telegramBotToken))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataString))
	hmacString := hex.EncodeToString(h.Sum(nil))

	return hmacString == hashValue
}

// Генерация JWT
func generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// Хэндлер логина через Telegram с JWT
func loginTelegramHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if checkResponse(data) {
		// Генерируем JWT
		token, err := generateJWT(data["id"])
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(map[string]string{
			"id":         data["id"],
			"first_name": data["first_name"],
			"last_name":  data["last_name"],
			"username":   data["username"],
			"photo_url":  data["photo_url"],
			"auth_date":  data["auth_date"],
			"token":      token, // сюда помещаем JWT
		})
		w.Write(response)
	} else {
		http.Error(w, "Authorization failed", http.StatusUnauthorized)
	}
}

func main() {
	http.HandleFunc("/login/telegram", withCORS(loginTelegramHandler))

	fmt.Println("✅ Server running on http://0.0.0.0:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
