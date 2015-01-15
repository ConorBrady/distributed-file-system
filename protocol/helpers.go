package protocol

import (
	"fmt"
	"log"
	"bytes"
	"runtime/debug"
)

// Yes I rewrote a reader wrapper around channels, this was a mistake, one I have
// learned a lesson from.

func read(input <-chan byte, delim byte) string {
	token := make([]byte,0)
	for msgByte := range input {
		if(msgByte==delim) {
			break
		}
		token = append(token,msgByte)
	}
	token = bytes.Trim(token,"\x00")
	return string(token)
}

func readByteCount(input <-chan byte, count int) []byte {
	token := make([]byte,0)
	i := 0
	for msgByte := range input {
		token = append(token,msgByte)
		i += 1
		if(i==count) {
			break
		}
	}
	return token
}

func readLine(input <-chan byte) string {
	line := read(input,'\n')
	log.Println(line)
	return line
}

func send(output chan<- byte, text string) {
	for _, b := range []byte(text) {
		output <- b
	}
}

func sendLine(output chan<- byte, text string) {
	log.Println(text)
	send(output,text+"\n")
}

func syncToToken(input <-chan byte, token string) {
	for _, r := range []byte(token) {
		for b := range input {
			if b == r {
				break
			}
		}
	}
}

func respondError(errorCode ErrorCode, output chan<- byte) {
	sendLine(output,fmt.Sprintf("ERROR_CODE: %d",errorCode))
	debug.PrintStack()
	switch errorCode {
		case ERROR_MALFORMED_REQUEST:
			sendLine(output,"ERROR_DESCRIPTION: MALFORMED REQUEST")
		case ERROR_UNKNOWN_PROTOCOL:
			sendLine(output,"ERROR_DESCRIPTION: PROTOCOL UNKNOWN")
		case ERROR_ROOM_NOT_FOUND:
			sendLine(output,"ERROR_DESCRIPTION: ROOM DOES NOT EXIST")
		case ERROR_CLIENT_NAME_MISMATCH:
			sendLine(output,"ERROR_DESCRIPTION: JOIN ID - CLIENT NAME MISMATCH")
		case ERROR_JOIN_ID_NOT_FOUND:
			sendLine(output,"ERROR_DESCRIPTION: JOIN ID NOT FOUND")
		case ERROR_FILE_NOT_FOUND:
			sendLine(output,"ERROR_DESCRIPTION: FILE NOT FOUND")
		case ERROR_USER_NOT_FOUND:
			sendLine(output,"ERROR_DESCRIPTION: USER NOT FOUND")
	}
}
