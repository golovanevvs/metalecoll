package agent

import (
	"crypto/sha256"
	"encoding/hex"
)

// calchash производит хэширование. Возвращает хэш.
func calcHash(data []byte, key string) string {
	h := sha256.New()
	h.Write(data)
	h.Write([]byte(key))
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)
}
