package chat

type ClientAdd struct {
	client *ChatClient
	response chan<- int
}

type ChatRoom struct {
	name string
	clients []*ChatClient
	queue chan ClientAdd
}

func MakeChatRoom(roomName string) *ChatRoom {
	r := &ChatRoom{
		roomName,
		make([]*ChatClient,0),
		make(chan ClientAdd,5),
	}
	go func() {
		for {
			clientAdd := <- r.queue
			r.clients = append(r.clients,clientAdd.client)
			clientAdd.response <- len(r.clients)-1
		}
	}()
	return r
}

func (r *ChatRoom)Name() string {
	return r.name
}

func (r *ChatRoom)Clients() []*ChatClient {
	return r.clients
}

func (r *ChatRoom)Add(client *ChatClient) <-chan int {
	channel := make(chan int)
	r.queue <- ClientAdd{
		client,
		channel,
	}
	return channel
}

func (r *ChatRoom)ClientForJoinId(joinId int) (*ChatClient, bool) {
	if len(r.clients) > joinId && r.clients[joinId].Valid() {
		return r.clients[joinId], true
	} else {
		return &ChatClient{}, false
	}
}
