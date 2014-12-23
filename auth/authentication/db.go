package authentication

import (
	"os"
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"log"
	)

func dbConnect() *sqlite3.Conn {

	db, _ := sqlite3.Open(os.Getenv("GOPATH")+"/src/distributed-file-system/auth/authentication/auth.sqlite")

	userCheck := 		"CREATE TABLE IF NOT EXISTS users(" +
							"user_id 	INTEGER PRIMARY KEY, "+
							"username 	VARCHAR(255) NOT NULL UNIQUE, "+
							"password 	VARCHAR(255) NOT NULL)"

	if err := db.Exec(userCheck); err != nil {
		log.Fatal(err.Error())
	}

	sessionKeyCheck :=	"CREATE TABLE IF NOT EXISTS session_keys(" +
							"user_id 	INTEGER NOT NULL, " +
							"key 		BLOB NOT NULL UNIQUE, " +
							"expiry 	VARCHAR(40) NOT NULL, " +
							"FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE )"

	if err := db.Exec(sessionKeyCheck); err != nil {
		log.Fatal(err.Error())
	}

	return db
}
