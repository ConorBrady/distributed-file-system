package auth

import (
	"os"
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"crypto/sha256"
	)

type User struct {
	userId int
	username string
	password string
}

func GetUser(username string) *User{

	db, _ := sqlite3.Open(os.Getenv("GOPATH")+"/src/distributed-file-system/auth/auth.sqlite")

	args := sqlite3.NamedArgs{
		"$username": username,
	}

	query, qErr := db.Query("select user_id, password from users where username=$username", args)

	if qErr != nil {
		return nil
	}

	var password string
	var userId int

	query.Scan(&userId,&password)
	query.Close()

	return &User{
		userId,
		username,
		password,
	}
}

func (u* User)getAESKey() []byte {
	key := sha256.Sum256([]byte(u.password))
	return key[:]
}
