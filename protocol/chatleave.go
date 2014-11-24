package protocol

import(
	"regexp"
	"../chat"
	"fmt"
	"strconv"
	)

type ChatLeaveProtocol struct {
	chat *chat.Chat
	queue chan *Exchange
}

func MakeChatLeaveProtocol(chat *chat.Chat, threadCount int) *ChatLeaveProtocol {
	p := &ChatLeaveProtocol{
		chat,
		make(chan *Exchange),
	}
	for i := 0; i < threadCount; i++ {
		go p.runLoop()
	}
	return p
}

func (p *ChatLeaveProtocol)Identifier() string {
	return "LEAVE_CHATROOM"
}

func (p *ChatLeaveProtocol)Handle(request <-chan byte, response chan<- byte) <-chan int {
	done := make(chan int)
	p.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (p *ChatLeaveProtocol)runLoop() {
	for {
		rr := <- p.queue

		// Line 1 "LEAVE_CHATROOM:"
		r1, _ := regexp.Compile("^\\s*(\\d*)\\s*$")
		matches1 := r1.FindStringSubmatch(readLine(rr.request))
		chatRoomRef, _ := strconv.Atoi(matches1[1])

		// Line 2 "JOIN_ID:"
		r2, _ := regexp.Compile("^JOIN_ID:\\s*(\\d*)\\s*$")
		matches2 := r2.FindStringSubmatch(readLine(rr.request))
		joinId, _ := strconv.Atoi(matches2[1])

		// Line 3 "CLIENT_NAME:"
		r3, _ := regexp.Compile("^CLIENT_NAME:\\s*(.*)\\s*$")
		clientName := r3.FindStringSubmatch(readLine(rr.request))[1]

		p.chat.LeaveRoom(chatRoomRef,joinId,clientName)

		response := fmt.Sprintf("LEFT_CHATROOM: %d\nJOIN_ID: %d\n",chatRoomRef,joinId)
		for _, b := range []byte(response) {
			rr.response <- b
		}
		rr.done <- 1
	}
}
