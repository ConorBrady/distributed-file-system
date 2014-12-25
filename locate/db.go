package locate

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"log"
	)

func dbConnect() *sqlite3.Conn {

	db, _ := sqlite3.Open("locate.sqlite")

	fileServers := 		"CREATE TABLE IF NOT EXISTS file_servers(" +
							"uuid 		VARCHAR(255) PRIMARY KEY, "+
							"address 	VARCHAR(255) NOT NULL )"

	if err := db.Exec(fileServers); err != nil {
		log.Fatal(err.Error())
	}

	files :=	"CREATE TABLE IF NOT EXISTS files(" +
					"name 				VARCHAR(255) PRIMARY KEY, " +
					"file_server_uuid 	VARCHAR(255) NOT NULL, "+
					"FOREIGN KEY(file_server_uuid) REFERENCES file_server(uuid) ON DELETE CASCADE )"

	if err := db.Exec(files); err != nil {
		log.Fatal(err.Error())
	}

	return db
}
