package service

import (
	"log"
	"os"
	"time"
	"regexp"
	"github.com/conorbrady/distributed-file-system/auth/crypto"
	"encoding/base64"
	"strconv"
	)

type SessionKey struct {
	key []byte
	username string
	expiry time.Time
}

func DecryptSessionKey(encryptedPackage []byte, keyPath string) *SessionKey {

	keyFile, _ := os.Open(keyPath)

	privateKey := make([]byte,32)

	keyFile.Read(privateKey)

	keyFile.Close()

	rgx, _ := regexp.Compile("\\s*SESSION_KEY:\\s*(\\S+)\\s+USERNAME:\\s*(\\w+)\\s+EXPIRES_AT:\\s*([\\d-]+T[\\d:\\.Z\\+]+)\\s*")
	matches := rgx.FindStringSubmatch(crypto.DecryptToString(encryptedPackage,privateKey))
	if len(matches) < 4 {
		log.Print(strconv.Itoa(len(matches))+" matches found")
		return nil
	}

	expiry, _ := time.Parse(time.RFC3339, matches[3])

	key, keyErr := base64.StdEncoding.DecodeString(matches[1])

	if keyErr != nil {
		return nil
	}
	return &SessionKey{
		key,
		matches[2],
		expiry,
	}

}

func (s* SessionKey)Key() []byte {
	return s.key
}

func (s* SessionKey)Username() string {
	return s.username
}
