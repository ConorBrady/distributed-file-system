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

type ASSessionKey struct {
	key []byte
	user User
	expiry time.Time
}

func GetASSessionKey(user User) *ASSessionKey {

	db, _ := sqlite3.Open(os.Getenv("GOPATH")+"/src/distributed-file-system/auth/auth.sqlite")

	nowBytes, _ := time.Now().MarshalText()

	queryArgs := sqlite3.NamedArgs{
		"$user_id": user.userId,
		"$now" : string(nowBytes),
	}

	query, qErr := db.Query("select key, expiry from session_keys where user_id = $user_id and expiry > $now", queryArgs)

	if qErr == nil {

		var key []byte
		var expiryStr string
		query.Scan(&key,&expiryStr)
		query.Close()

		expiry, _ := time.Parse(time.RFC3339, expiryStr)

		return &ASSessionKey{
			key,
			user,
			expiry,
		}

	} else {

		keyDuration, _ := time.ParseDuration("12h")
		expiry := time.Now().Add(keyDuration)
		expiryBytes, _ := expiry.MarshalText()

		key := make([]byte, 32)
		rand.Read(key)

		insertArgs := sqlite3.NamedArgs{
			"$user_id"  : user.userId,
			"$key"		: key,
			"$expiry"	: string(expiryBytes),
		}

		createErr := db.Exec( "insert into session_keys ( user_id, key, expiry ) values ( $user_id, $key, $expiry )", insertArgs )

		if createErr == nil {
			return &ASSessionKey{
				key,
				user,
				expiry,
			}
		} else {
			fmt.Println(createErr.Error())
			return nil
		}
	}
}

func (s* ASSessionKey)EncryptedKey() []byte {

	encKey := make([]byte,32)

	cipherBlock, _ := aes.NewCipher(s.user.getAESKey())
	cipherBlock.Encrypt(encKey[0:16], s.key[0:16])
	cipherBlock.Encrypt(encKey[16:32],s.key[16:32])

	return encKey
}

func (s* ASSessionKey)MarshalAndEncrypt() []byte {

	expiryBytes, _ := s.expiry.MarshalText()

	packageStr := 	"SESSION_KEY: "+base64.StdEncoding.EncodeToString(s.key)+"\n"+
				  	"USERNAME: "+s.user.username+"\n"+
					"EXPIRES_AT: "+string(expiryBytes)+"\n"

	packageBytes := make([]byte,int(math.Ceil(float64(len(packageStr))/16.0))*16)

	for i, b := range packageStr {
		packageBytes[i] = byte(b)
	}

	keyFile, _ := os.Open(os.Getenv("GOPATH")+"/src/distributed-file-system/auth/private.key")

	privateKey := make([]byte,32)
	keyFile.Read(privateKey)

	keyFile.Close()

	cipherBlock, _ := aes.NewCipher(privateKey)

	encryptedPackage := make([]byte,len(packageBytes))

	for i := 0; i < len(packageBytes); i += aes.BlockSize {
		cipherBlock.Encrypt(encryptedPackage[i:i+aes.BlockSize], packageBytes[i:i+aes.BlockSize])
	}

	return encryptedPackage
}
