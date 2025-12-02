package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

func (s *serv) VerifyTelegramData(data map[string]string) bool {
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

	secretKey := sha256.Sum256([]byte(s.telegramBotToken))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataString))
	hmacString := hex.EncodeToString(h.Sum(nil))

	return hmacString == hashValue
}
