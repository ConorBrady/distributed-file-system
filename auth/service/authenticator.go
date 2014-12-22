package service

import (
	"strconv"
	"regexp"
	"log"
	"encoding/base64"
	"distributed-file-system/auth/crypto"
	)

type Authenticator struct {
	username string
	timestamp int
}

func DecryptAuthenticator(encryptedPackage []byte, key []byte) *Authenticator {

	rgx, _ := regexp.Compile("\\s*USERNAME:\\s*(\\w+)\\s+TIMESTAMP:\\s*(\\d+)\\s*")
	matches := rgx.FindStringSubmatch(crypto.DecryptToString(encryptedPackage,key))
	if len(matches) < 3 {
		log.Print("Decryption failed... Provided Key: "+base64.StdEncoding.EncodeToString(key))
		return nil
	}

	timestamp, err := strconv.Atoi(matches[2])

	if err != nil {
		log.Print("Invalid timestamp: " + matches[2])
		return nil
	}

	return &Authenticator{
		matches[1],
		timestamp,
	}
}

func (a* Authenticator)Username() string {
	return a.username
}

func (a* Authenticator)MakeResponse(key []byte) []byte {

	return crypto.EncryptString(strconv.Itoa(a.timestamp+1),key)
}
