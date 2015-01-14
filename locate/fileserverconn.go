package locate

import (
	"bufio"
	"net"
	"regexp"
	)
type FileServerConn struct {
	reader *bufio.Reader
	writer *net.TCPConn
}

func FSConnect(address string) *FileServerConn {

	// File Server Communication

	fileServAddr, fsAddrErr := net.ResolveTCPAddr("tcp4",address)

	if fsAddrErr != nil {
		return nil
	}

	fileServConn, fsConnErr := net.DialTCP("tcp4",nil,fileServAddr)

	if fsConnErr != nil {
		return nil
	}

	fileServConnReader := bufio.NewReader(fileServConn)

	return &FileServerConn{
		fileServConnReader,
		fileServConn,
	}
}

func (c *FileServerConn)getUUID() *string {

	c.writer.Write([]byte("HELO:\n"))

	line, _ := c.reader.ReadString('\n')
	r, _ := regexp.Compile("\\AHELO:\\s*\\S*\\s*\\z")
	if !r.MatchString(line) {
		return nil
	}

	line, _ = c.reader.ReadString('\n')
	r, _ = regexp.Compile("\\AIP:\\s*\\S*\\s*\\z")
	if !r.MatchString(line) {
		return nil
	}

	line, _ = c.reader.ReadString('\n')
	r, _ = regexp.Compile("\\APort:\\s*\\S*\\s*\\z")
	if !r.MatchString(line) {
		return nil
	}

	line, _ = c.reader.ReadString('\n')
	r, _ = regexp.Compile("\\AStudentID:\\s*\\S*\\s*\\z")
	if !r.MatchString(line) {
		return nil
	}

	line, _ = c.reader.ReadString('\n')
	r, _ = regexp.Compile("\\AUUID:\\s*(\\S+)\\s*\\z")
	matches := r.FindStringSubmatch(line)
	if len(matches) < 2 {
		return nil
	}

	line, _ = c.reader.ReadString('\n')
	r, _ = regexp.Compile("\\AMODE:\\s*FS\\s*\\z")
	if !r.MatchString(line) {
		return nil
	}

	return &matches[1]
}
