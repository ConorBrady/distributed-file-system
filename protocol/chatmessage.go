package protocol

import(
	"regexp"
	"../chat"
	"strconv"
	)

type ChatMessageProtocol struct {
	chat *chat.Chat
	queue chan *Exchange
}

func MakeChatMessageProtocol(chat *chat.Chat, threadCount int) *ChatMessageProtocol {
	p := &ChatMessageProtocol{
		chat,
		make(chan *Exchange),
	}
	for i := 0; i < threadCount; i++ {
		go p.runLoop()
	}
	return p
}

func (p *ChatMessageProtocol)Identifier() string {
	return "CHAT"
}

func (p *ChatMessageProtocol)Handle(request <-chan byte, response chan<- byte) <-chan int {
	done := make(chan int)
	p.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (p *ChatMessageProtocol)runLoop() {
	for {
		rr := <- p.queue

		// Line 1 "CHAT:"
		r1, _ := regexp.Compile("\\A\\s*(\\d*)\\s*\\z")
		matches1 := r1.FindStringSubmatch(readLine(rr.request))
		chatRoomRef, _ := strconv.Atoi(matches1[1])

		// Line 2 "JOIN_ID:"
		r2, _ := regexp.Compile("\\AJOIN_ID:\\s*(\\d*)\\s*\\z")
		matches2 := r2.FindStringSubmatch(readLine(rr.request))
		joinId, _ := strconv.Atoi(matches2[1])

		// Line 3 "CLIENT_NAME:"
		r3, _ := regexp.Compile("\\ACLIENT_NAME:\\s*(.*)\\s*\\z")
		matches3 := r3.FindStringSubmatch(readLine(rr.request))
		clientName := matches3[1]

		// Line 4 "MESSAGE:"
		r4, _ := regexp.Compile("\\AMESSAGE:\\s*(.*)\\s*\\z")
		matches4 := r4.FindStringSubmatch(readLine(rr.request))
		message := matches4[1]

		p.chat.Rooms()[chatRoomRef].Send(message,clientName,joinId)
		rr.done <- 1
	}
}
