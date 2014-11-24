package protocol

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
