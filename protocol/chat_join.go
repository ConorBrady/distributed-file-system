package protocol

import(
	"regexp"
	"github.com/conorbrady/distributed-file-system/chat"
	"fmt"
	)

type ChatJoinProtocol struct {
	chat *chat.Chat
	queue chan *Exchange
}

func MakeChatJoinProtocol(chat *chat.Chat, threadCount int) *ChatJoinProtocol {
	p := &ChatJoinProtocol{
		chat,
		make(chan *Exchange, threadCount),
	}
	for i := 0; i < threadCount; i++ {
		go p.runLoop()
	}
	return p
}

func (p *ChatJoinProtocol)Identifier() string {
	return "JOIN_CHATROOM"
}

func (p *ChatJoinProtocol)Handle(request <-chan byte, response chan<- byte) <-chan StatusCode{
	done := make(chan StatusCode, 1)
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
		r1, _ := regexp.Compile("\\A\\s*(\\w+.*)\\s*\\z")
		matches1 := r1.FindStringSubmatch(readLine(rr.request))
		if len(matches1) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}
		chatRoomName := string(matches1[1])

		// Line 2 "CLIENT_IP: 0"
		r2, _ := regexp.Compile("\\ACLIENT_IP:\\s*0\\s*\\z")
		if !r2.MatchString(readLine(rr.request)) {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		// Line 3 "PORT: 0"
		r3, _ := regexp.Compile("\\APORT:\\s*0\\s*\\z")
		if !r3.MatchString(readLine(rr.request)) {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		// Line 4 "CLIENT_NAME:"
		r4, _ := regexp.Compile("\\ACLIENT_NAME:\\s*(\\w+.*)\\s*\\z")
		matches4 := r4.FindStringSubmatch(readLine(rr.request))
		if len(matches4) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}
		clientName := string(matches4[1])

		r := <- p.chat.JoinRoom(chatRoomName,chat.MakeChatClient(clientName,rr.response))
		roomRef := r.RoomId
		joinId := <- r.JoinId

		sendLine(rr.response,fmt.Sprintf("JOINED_CHATROOM: %s",chatRoomName))
		sendLine(rr.response,"SERVER_IP: 0")
		sendLine(rr.response,"PORT: 0")
		sendLine(rr.response,fmt.Sprintf("ROOM_REF: %d",roomRef))
		sendLine(rr.response,fmt.Sprintf("JOIN_ID: %d",joinId))

		rr.done <- STATUS_SUCCESS_CONTINUE
	}
}
