package protocol

import(
	"regexp"
	"../chat"
)

type DisconnectProtocol struct {
	chat *chat.Chat
	queue chan *Exchange
}

func MakeDisconnectProtocol(chat *chat.Chat, threadCount int) *DisconnectProtocol {
	p := &DisconnectProtocol{
		chat,
		make(chan *Exchange),
	}
	for i := 0; i < threadCount; i++ {
		go p.runLoop()
	}
	return p
}

func (p *DisconnectProtocol)Identifier() string {
	return "DISCONNECT"
}

func (p *DisconnectProtocol)Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
	done := make(chan StatusCode)
	p.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (p *DisconnectProtocol)runLoop() {
	for {
		rr := <- p.queue

		// Line 1 "DISCONNECT:"
		r1, _ := regexp.Compile("\\A\\s*0\\s*\\z")
		if !r1.MatchString(readLine(rr.request)) {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		// Line 2 "PORT:"
		r2, _ := regexp.Compile("\\APORT:\\s*0\\s*\\z")
		if !r2.MatchString(readLine(rr.request)) {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		// Line 3 "CLIENT_NAME:"
		r3, _ := regexp.Compile("\\ACLIENT_NAME:\\s*(\\w.*)\\s*\\z")
		if !r3.MatchString(readLine(rr.request)) {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}
		for _, room := range p.chat.Rooms() {
			for _, client := range room.Clients() {
				if client.Channel() == rr.response {
					client.Invalidate()
				}
			}
		}
		close(rr.response)
		rr.done <- STATUS_SUCCESS_DISCONNECT
	}
}
