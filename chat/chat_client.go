package chat

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

func (c *ChatClient)Channel() chan<- byte {
	return c.connection
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
