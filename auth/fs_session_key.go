package auth

import (
	"os"
	"fmt"
	"time"
	"crypto/rand"
	"math"
	"crypto/aes"
	"encoding/base64"
	"code.google.com/p/go-sqlite/go1/sqlite3"
	)

type FSSessionKey struct {
	key []byte
	username string
	expiry time.Time
}

func GenFromTicket(ticket string) *FSSessionKey {

	keyFile, _ := os.Open(os.Getenv("GOPATH")+"/src/distributed-file-system/auth/private.key")

	privateKey := make([]byte,32)
	keyFile.Read(privateKey)

	keyFile.Close()

	cipherBlock, _ := aes.NewCipher(privateKey)

	encryptedPackage := base64.StdEncoding.DecodeString(ticket)
	packageBytes := make([]byte,len(encryptedPackage))

	for i := 0; i < len(packageBytes); i += aes.BlockSize {
		cipherBlock.Decrypt(packageBytes[i:i+aes.BlockSize], encryptedPackage[i:i+aes.BlockSize] )
	}

	rgx, _ := regexp.Compile("\\A\\s*SESSION_KEY:\\s*(\\w+=*)\\s+USERNAME:\\s*(\\w+)\\s+EXPIRES_AT:\\s*([\\d-]+T[\\d:\\+]+)\\s*\\z")
	matches := rgx.FindStringSubmatch(string(packageBytes))
	if len(matches) < 4 {
		return nil
	}

	expiry, _ := time.Parse(time.RFC3339, matches[3])

	return &FSSessionKey{
		base64.StdEncoding.DecodeString(matches[1]),
		matches[2],
		expiry,
	}

}
