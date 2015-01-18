package file

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"log"
)

const MAX_BLOCK_SIZE = 4096

type Block struct {
	filename string
	index int
	hash string
	size int
}

func (b *Block)Hash() string {
	return b.hash
}

func GetBlock(filename string, index int) (*Block, error) {
	return getBlock(filename, index)
}

func getBlock(filename string, index int) (*Block, error) {

	args := sqlite3.NamedArgs{
		"$filename"	: filename,
		"$index"	: index,
	}

	query, qErr := dbConnect().Query("select hash, size from blocks where filename = $filename and block_index = $index",args)

	if qErr != nil {
		log.Println("Block lookup failed with error: "+qErr.Error())
		return nil, qErr
	}

	var hash string
	var size int

	query.Scan(&hash,&size)
	query.Close()

	return &Block{
		filename,
		index,
		hash,
		size,
	}, nil
}

func (b* Block)setHash(hash string) error {

	args := sqlite3.NamedArgs{
		"$oldHash"	: b.hash,
		"$newHash"	: hash,
	}

	b.hash = hash

	return dbConnect().Exec("update blocks set hash = $newHash where hash = $oldHash", args)
}

func (b* Block)setSize(size int) error {

	args := sqlite3.NamedArgs{
		"$hash"	: b.hash,
		"$size" : size,
	}

	b.size = size

	return dbConnect().Exec("update blocks set size = $size where hash = $hash", args)
}

func createBlock(filename string, index int, hash string, size int) error {

	args := sqlite3.NamedArgs{
		"$filename" : filename,
		"$index"	: index,
		"$hash"		: hash,
		"$size"		: size,
	}

	return dbConnect().Exec("insert into blocks (filename, block_index, hash, size) values ($filename, $index, $hash, $size) ", args)
}
