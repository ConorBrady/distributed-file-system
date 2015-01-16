package file

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"log"
	)

func dbConnect() *sqlite3.Conn {

	db, _ := sqlite3.Open("file.sqlite")

	blockCheck := 	"CREATE TABLE IF NOT EXISTS blocks(" +
						"filename 		VARCHAR(255) NOT NULL, "+
						"block_index	INTEGER, "+
						"hash 			VARCHAR(255) NOT NULL, "+
						"size			INTEGER NOT NULL, "+
						"PRIMARY KEY( filename, block_index ))"

	if err := db.Exec(blockCheck); err != nil {
		log.Println("Database initialization error")
		log.Fatal(err.Error())
	}

	return db
}
