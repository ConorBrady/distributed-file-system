package chat

import(
	"strconv"
	)
type ChatClient struct {
	name string
	connection chan<- byte
	valid bool
}

func MakeChatClient(name string, connection chan<- byte) *ChatClient {
	return &ChatClient{
		name,
		connection,
		true,
	}
}

func (c *ChatClient)Send(message string, from string, roomRef int) {
	message = "CHAT: "+strconv.Itoa(roomRef)+"\nCLIENT_NAME: "+from+"\nMESSAGE: "+message+"\n"
	for _, b := range []byte(message) {
		c.connection <- b
	}
}

func (c *ChatClient)Name() string {
	return c.name
}

func (c *ChatClient)Valid() bool {
	return c.valid
}

func (c *ChatClient)Invalidate() {
	c.valid = false
}
