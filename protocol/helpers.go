package protocol

import (
	"fmt"
)

func readLine(input <-chan byte) string {
	line := make([]byte,0)
	for msgByte := range input {
		if(msgByte=='\n') {
			break
		}
		line = append(line,msgByte)
	}
	return string(line)
}

func sendLine(output chan<- byte, text string) {
	for _, b := range []byte(text+"\n") {
		output <- b
	}
}

func error(errorCode ErrorCode, output chan<- byte) {
	sendLine(output,fmt.Sprintf("ERROR_CODE: %d",errorCode))
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
	}
}
