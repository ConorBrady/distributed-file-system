package file

import (

	"os"
	"errors"
	"log"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
)

func ReadData(filename string, index int) ([]byte, error) {

	os.Mkdir("storage",0777)

	if block, err := getBlock(filename,index); err == nil {
		data, _ := ioutil.ReadFile("storage/"+block.hash)
		return data[ :block.size], nil
	} else {
		return nil, err
	}
}

func WriteData(filename string, start int, data []byte) error {

	os.Mkdir("storage",0777)

	index := start/MAX_BLOCK_SIZE

	// If its not the first and it has no proceeding full block it is invalid
	if index != 0 {
		oldBlock, _ := getBlock(filename, index-1)
	 	if oldBlock == nil || oldBlock.size != MAX_BLOCK_SIZE {
			log.Println("WARNING: Attempted to print non-contiguous block")
			return errors.New("Cannot create non-contiguous block")
		}
	}

	// Loop here for writes across blocks
	{
		buffer := make([]byte,MAX_BLOCK_SIZE)

		oldBlock, blErr := getBlock(filename, index)

		if blErr == nil {

			file, fiErr := os.Open("storage/"+oldBlock.hash)
			if fiErr != nil {
				log.Println("Error opening old block: "+fiErr.Error())
				return fiErr
			}
			file.Read(buffer)
			file.Close()
		}

		offset := start%MAX_BLOCK_SIZE
		writeEnd := offset+len(data)

		for i := 0; i < len(data) && i+offset < len(buffer); i++ {
			buffer[i+offset] = data[i]
		}

		hasher := sha1.New()
		hasher.Write(buffer)
		hash := hex.EncodeToString(hasher.Sum(nil))

		if oldBlock == nil || oldBlock.hash != hash {

			err := ioutil.WriteFile("storage/"+hash, buffer, 0777)

			if err != nil {
				log.Println("Error writing new block: "+err.Error())
				return err
			}
		}

		if oldBlock == nil {

			log.Println("Creating new block")
			if err := createBlock(filename,index,hash,writeEnd); err != nil {
				log.Println("An error occured creating block: "+err.Error())
				return err
			}

		} else if oldBlock.hash != hash {

			log.Println("Updating database for new block")
			os.Remove("storage/"+oldBlock.hash)

			oldBlock.setHash(hash)

			if oldBlock.size < writeEnd {
				oldBlock.setSize(writeEnd)
			}
		}

	}

	return nil
}
