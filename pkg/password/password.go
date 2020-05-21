package pass

import (
	"bytes"
	"crypto/rand"

	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/argon2"
)

const (
	lenSalt = 8
)

func HashPassword(salt, password []byte) []byte {
	hashPass := argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
	return append(salt, hashPass...)
}

func HashPasswordGenSalt(password []byte) []byte {
	salt := make([]byte, lenSalt)
	_, err := rand.Read(salt)
	if err != nil {
		log.Error(err)
	}
	return HashPassword(salt, password)
}

func CheckPassword(pass, realHashPass []byte) bool {
	var salt []byte
	salt = append(salt, realHashPass[0:lenSalt]...)
	return bytes.Equal(HashPassword(salt, pass), realHashPass)
}
