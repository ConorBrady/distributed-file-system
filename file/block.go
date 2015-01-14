package file

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"fmt"
	"log"
)

const BlockSize = 4096

type Block struct {
	filename string
	index int
	hash string
}

func (b *Block)Hash() string {
	return b.hash
}

func GetBlock(filename string, index int) (*Block, error) {
	return getBlock(filename, index)
}

func getBlock(filename string, index int) (*Block, error) {

	log.Println(fmt.Sprintf("Looking up block for %s and index %d", filename, index ))
	args := sqlite3.NamedArgs{
		"$filename"	: filename,
		"$index"	: index,
	}

	query, qErr := dbConnect().Query("select hash from blocks where filename = $filename and block_index = $index",args)

	if qErr != nil {
		log.Println("Block lookup failed with error: "+qErr.Error())
		return nil, qErr
	}

	var hash string

	query.Scan(&hash)
	fmt.Println(hash)
	query.Close()

	return &Block{
		filename,
		index,
		hash,
	}, nil
}

func (b* Block)setHash(hash string) error {

	args := sqlite3.NamedArgs{
		"$oldHash"	: b.hash,
		"$newHash"	: hash,
	}

	return dbConnect().Exec("update blocks set hash = $newHash where hash = $oldHash", args)
}

func createBlock(filename string, index int, hash string) error {

	args := sqlite3.NamedArgs{
		"$filename" : filename,
		"$index"	: index,
		"$hash"		: hash,
	}

	return dbConnect().Exec("insert into blocks (filename, block_index, hash) values ($filename, $index, $hash) ", args)
}
