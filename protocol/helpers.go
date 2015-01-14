package protocol

import (
	"fmt"
	"io"
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
	return read(input,'\n')
}

func send(output chan<- byte, text string) {
	for _, b := range []byte(text) {
		output <- b
	}
}

func sendLine(output chan<- byte, text string) {
	send(output,text+"\n")
}

type ChannelWriter struct {
	channel chan<- byte
}

func MakeChannelWriter(channel chan<- byte) *ChannelWriter {
	c := &ChannelWriter{
		channel,
	}
	return c
}

func (c* ChannelWriter)Write(p []byte) (n int, err error) {
	fmt.Println("Called writer")
	for _, b := range p {
		fmt.Println("Sending "+string(b))
		c.channel <- b
	}
	return len(p), nil
}

type ChannelReader struct {
	channel <-chan byte
}

func MakeChannelReader(channel <-chan byte) *ChannelReader {
	c := &ChannelReader{
		channel,
	}
	return c
}

func (c* ChannelReader)Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		b := <- c.channel
		if b == ' ' || b == '\n' {
			fmt.Println("End of file")
			return i-1, io.EOF
		} else {
			p[i] = b
		}
	}
	return len(p), nil
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
