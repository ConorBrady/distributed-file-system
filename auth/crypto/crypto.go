package crypto

import (
	"crypto/aes"
	"bytes"
	"math"
	)

func EncryptString(line string, key []byte) []byte {

	bytes := make([]byte,int(math.Ceil(float64(len(line))/aes.BlockSize))*aes.BlockSize)

	for i, b := range line {
		bytes[i] = byte(b)
	}

	return EncryptBytes(bytes,key)
}

func EncryptBytes(data []byte, key []byte) []byte {

	cipherBlock, _ := aes.NewCipher(key)

	enc := make([]byte,len(data))

	for i := 0; i < len(data); i += aes.BlockSize {
		cipherBlock.Encrypt(enc[i:i+aes.BlockSize], data[i:i+aes.BlockSize])
	}

	return enc
}

func DecryptToBytes(enc []byte, key []byte) []byte {
	data := make([]byte,len(enc))

	cipherBlock, _ := aes.NewCipher(key)

	for i := 0; i < len(enc); i += aes.BlockSize {
		cipherBlock.Decrypt(data[i:i+aes.BlockSize], enc[i:i+aes.BlockSize])
	}

	return bytes.Trim(data,"\x00")
}

func DecryptToString(enc []byte, key []byte) string {
	return string(DecryptToBytes(enc,key))
}
