package protocol

import(
	"regexp"
	"../chat"
	"fmt"
	)

type ChatJoinProtocol struct {
	chat *chat.Chat
	queue chan *Exchange
}

func MakeChatJoinProtocol(chat *chat.Chat, threadCount int) *ChatJoinProtocol {
	p := &ChatJoinProtocol{
		chat,
		make(chan *Exchange),
	}
	for i := 0; i < threadCount; i++ {
		go p.runLoop()
	}
	return p
}

func (p *ChatJoinProtocol)Identifier() string {
	return "JOIN_CHATROOM"
}

func (p *ChatJoinProtocol)Handle(request <-chan byte, response chan<- byte) <-chan int{
	done := make(chan int)
	p.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (p *ChatJoinProtocol)runLoop() {
	for {
		rr := <- p.queue

		// Line 1 "JOIN_CHATROOM:"
		r1, _ := regexp.Compile("\\A\\s*(.*)\\s*\\z")
		matches1 := r1.FindStringSubmatch(readLine(rr.request))
		chatRoomName := string(matches1[1])

		// Line 2 "CLIENT_IP: 0"
		// _, _ := regexp.Compile("\\ACLIENT_IP:\\s*(.*)\\s*\\z")
		readLine(rr.request)

		// Line 3 "PORT: 0"
		// _, _ := regexp.Compile("\\APORT:\\s*(.*)\\s*\\z")
		readLine(rr.request)

		// Line 4 "CLIENT_NAME:"
		r4, _ := regexp.Compile("\\ACLIENT_NAME:\\s*(.*)\\s*\\z")
		clientName := r4.FindStringSubmatch(readLine(rr.request))[1]

		r := <- p.chat.JoinRoom(chatRoomName,*chat.MakeChatClient(clientName,rr.response))
		roomRef := r.RoomId
		joinId := <- r.JoinId

		response := fmt.Sprintf("JOINED_CHATROOM: %s\nSERVER_IP: 0\nPORT: 0\nROOM_REF: %d\nJOIN_ID: %d\n",chatRoomName,roomRef,joinId)
		for _, b := range []byte(response) {
			rr.response <- b
		}
		rr.done <- 1
	}
}
