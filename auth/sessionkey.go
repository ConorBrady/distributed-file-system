package auth

import (
	"os"
	"fmt"
	"time"
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go-sqlite/go1/sqlite3"
	)

type SessionKey struct {
	key string
	user User
	expiry time.Time
}

func GetSessionKey(user User) *SessionKey {

	db, _ := sqlite3.Open(os.Getenv("GOPATH")+"/src/distributed-file-system/auth/auth.sqlite")

	nowBytes, _ := time.Now().MarshalText()

	queryArgs := sqlite3.NamedArgs{
		"$user_id": user.userId,
		"$now" : string(nowBytes),
	}

	query, qErr := db.Query("select key, expiry from session_keys where user_id = $user_id and expiry > $now", queryArgs)

	if qErr == nil {

		var key string
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

		key := uuid.New()

		insertArgs := sqlite3.NamedArgs{
			"$user_id"  : user.userId,
			"$key"		: key,
			"$expiry"	: string(expiryBytes),
		}

		createErr := db.Exec( "insert into session_keys ( user_id, key, expiry ) values ( $user_id, $key, $expiry )", insertArgs )

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

func (s* SessionKey)Key() string {
	return s.key
}
