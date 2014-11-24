package chat

type RoomAddResponse struct {
	RoomId int
	JoinId <-chan int
}

type RoomJoin struct {
	name string
	client ChatClient
	response chan<- RoomAddResponse
}

type Chat struct {
	rooms []ChatRoom
	queue chan RoomJoin
}

func MakeChat() *Chat{
	c := &Chat{
		make([]ChatRoom,0),
		make(chan RoomJoin),
	}
	go func(){
		for {
			roomAdd := <- c.queue
			roomId := c.roomIdForName(roomAdd.name)
			if roomId == -1 {
				c.rooms = append(c.rooms,*MakeChatRoom(roomAdd.name,len(c.rooms)))
				roomId = len(c.rooms)-1
			}
			roomAdd.response <- RoomAddResponse{
				roomId,
				c.rooms[roomId].Add(roomAdd.client),
			}
		}
	}()
	return c
}

func (c *Chat)roomIdForName(name string) int {
	for i, room := range c.rooms {
		if room.Name() == name {
			return i
		}
	}
	return -1
}

func (c *Chat)JoinRoom(roomName string, client ChatClient) <-chan RoomAddResponse {
	response := make(chan RoomAddResponse)
	c.queue <- RoomJoin{
		roomName,
		client,
		response,
	}
	return response
}

func (c *Chat)LeaveRoom(roomRef int, joinId int, clientName string) bool{
	if len(c.rooms) > roomRef {
		return c.rooms[roomRef].Remove(joinId,clientName)
	} else {
		return false
	}
}

func (c *Chat)Rooms() []ChatRoom {
	return c.rooms
}
