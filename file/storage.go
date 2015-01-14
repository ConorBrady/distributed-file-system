package file

import (

	"os"
	"errors"

	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
)

func ReadData(filename string, index int) ([]byte, error) {

	if block, err := getBlock(filename,index); err == nil {
		return ioutil.ReadFile("storage/"+block.hash)
	} else {
		return nil, err
	}
}

func WriteData(filename string, start int, data []byte) error {

	index := start/BlockSize

	// If its not the first and it has no proceeding block it is invalid
	if index != 0 {
		oldBlock, _ := getBlock(filename, index-1)
	 	if oldBlock == nil {
			return errors.New("Cannot create non-contiguous block")
		}
	}

	// Loop here for writes across blocks
	{
		buffer := make([]byte,BlockSize)

		oldBlock, blErr := getBlock(filename, index)

		if blErr == nil {

			file, fiErr := os.Open("storage/"+oldBlock.hash)
			if fiErr != nil {
				return fiErr
			}
			file.Read(buffer)
			file.Close()
		}

		offset := start%BlockSize

		for i := 0; i < len(data) && i+offset < len(buffer); i++ {
			buffer[i+offset] = data[i]
		}

		hasher := sha1.New()
		hasher.Write(buffer)
		hash := hex.EncodeToString(hasher.Sum(nil))

		if oldBlock == nil || oldBlock.hash != hash {

			err := ioutil.WriteFile("storage/"+hash, buffer, 0777)

			if err != nil {
				return err
			}
		}

		if oldBlock == nil {

			createBlock(filename,index,hash)

		} else if oldBlock.hash != hash {

			os.Remove("storage/"+oldBlock.hash)

			oldBlock.setHash(hash)
		}

	}

	return nil
}
