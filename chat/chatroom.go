package chat

type ClientAdd struct {
	client ChatClient
	response chan<- int
}

type ChatRoom struct {
	name string
	clients []ChatClient
	roomId int
	queue chan ClientAdd
}

func MakeChatRoom(roomName string, roomId int) *ChatRoom {
	r := &ChatRoom{
		roomName,
		make([]ChatClient,0),
		roomId,
		make(chan ClientAdd),
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

func (r *ChatRoom)Send(message string, from string, omit int) bool {

	if len(r.clients) > omit && r.clients[omit].Name() == from {

		for _, client := range r.clients {
			if client.Name() != from && client.Valid() {
				client.Send(message,from,r.roomId)
			}
		}
		return true
	} else {
		return false
	}
}

func (r *ChatRoom)Add(client ChatClient) <-chan int {
	channel := make(chan int)
	r.queue <- ClientAdd{
		client,
		channel,
	}
	return channel
}

func (r *ChatRoom)Remove(joinId int, name string) bool {
	if len(r.clients) > joinId && r.clients[joinId].Name() == name && r.clients[joinId].Valid() {
		r.clients[joinId].Invalidate()
		return true
	} else {
		return false
	}
}
