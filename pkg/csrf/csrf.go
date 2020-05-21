package csrf

import (
	"crypto/sha256"
	"encoding/base64"
	"io"

	"github.com/labstack/gommon/log"
)

const (
	CSRFheader = "X-CSRF-TOKEN"
	// 32 bytes
	csrfSalt = "KbWaoi5xtDC3GEfBa9ovQdzOzXsuVU9I"
)

func MakeToken(sid string) string {
	hasher := sha256.New()
	_, err := io.WriteString(hasher, csrfSalt+sid)
	if err != nil {
		log.Error(err)
	}
	token := base64.RawStdEncoding.EncodeToString(hasher.Sum(nil))
	return token
}

func ValidateToken(token string, sid string) bool {
	trueToken := MakeToken(sid)
	return token == trueToken
}
