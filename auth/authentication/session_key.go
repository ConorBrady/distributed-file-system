package authentication

import (
	"os"
	"fmt"
	"time"
	"crypto/rand"
	"distributed-file-system/auth/crypto"
	"encoding/base64"
	"code.google.com/p/go-sqlite/go1/sqlite3"
	)

type SessionKey struct {
	key []byte
	user User
	expiry time.Time
}

func GetSessionKey(user User) *SessionKey {

	nowBytes, _ := time.Now().MarshalText()

	queryArgs := sqlite3.NamedArgs{
		"$user_id": user.userId,
		"$now" : string(nowBytes),
	}

	query, qErr := dbConnect().Query("select key, expiry from session_keys where user_id = $user_id and expiry > $now", queryArgs)

	if qErr == nil {

		var key []byte
		var expiryStr string
		query.Scan(&key,&expiryStr)
		query.Close()

		expiry, _ := time.Parse(time.RFC3339, expiryStr)

		return &SessionKey{
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

		createErr := dbConnect().Exec( "insert into session_keys ( user_id, key, expiry ) values ( $user_id, $key, $expiry )", insertArgs )

		if createErr == nil {
			return &SessionKey{
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

func (s* SessionKey)EncryptedKey() []byte {

	return crypto.EncryptBytes(s.key,s.user.getAESKey())
}

func (s* SessionKey)MarshalAndEncrypt(keyPath string) []byte {

	expiryBytes, _ := s.expiry.MarshalText()

	packageStr := 	"SESSION_KEY: "+base64.StdEncoding.EncodeToString(s.key)+"\n"+
				  	"USERNAME: "+s.user.username+"\n"+
					"EXPIRES_AT: "+string(expiryBytes)+"\n"

	keyFile, _ := os.Open(keyPath)

	privateKey := make([]byte,32)
	keyFile.Read(privateKey)

	keyFile.Close()

	return crypto.EncryptString(packageStr,privateKey)
}
