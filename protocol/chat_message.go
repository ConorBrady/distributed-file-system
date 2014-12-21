package protocol

import(
	"regexp"
	"distributed-file-system/chat"
	"strconv"
	"fmt"
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

func (p *ChatMessageProtocol)Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
	done := make(chan StatusCode)
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
		r1, _ := regexp.Compile("\\A\\s*(\\d+)\\s*\\z")
		matches1 := r1.FindStringSubmatch(readLine(rr.request))
		if len(matches1) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}
		chatRoomRef, _ := strconv.Atoi(matches1[1])
		chatRoom, ok1 := p.chat.RoomForRef(chatRoomRef)
		if !ok1 {
			respondError(ERROR_ROOM_NOT_FOUND,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		// Line 2 "JOIN_ID:"
		r2, _ := regexp.Compile("\\AJOIN_ID:\\s*(\\d+)\\s*\\z")
		matches2 := r2.FindStringSubmatch(readLine(rr.request))
		if len(matches2) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}
		joinId, _ := strconv.Atoi(matches2[1])

		client, ok2 := chatRoom.ClientForJoinId(joinId)
		if !ok2 {
			respondError(ERROR_JOIN_ID_NOT_FOUND,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		// Line 3 "CLIENT_NAME:"
		r3, _ := regexp.Compile("\\ACLIENT_NAME:\\s*(\\w.*)\\s*\\z")
		matches3 := r3.FindStringSubmatch(readLine(rr.request))
		if len(matches3) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}
		clientName := matches3[1]
		if clientName != client.Name() {
			respondError(ERROR_CLIENT_NAME_MISMATCH,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		// Line 4 "MESSAGE:"
		r4, _ := regexp.Compile("\\AMESSAGE:\\s*(\\w.*)\\s*\\z")
		matches4 := r4.FindStringSubmatch(readLine(rr.request))
		if len(matches4) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}
		messageLines := []string{}

		for len(matches4) >= 2 {
			messageLines = append(messageLines,matches4[1])
			r4, _ := regexp.Compile("\\A\\s*(\\w.*)\\s*\\z")
			matches4 = r4.FindStringSubmatch(readLine(rr.request))
		}

		for _, recipient := range chatRoom.Clients() {
			if recipient != client && recipient.Valid() {
				sendLine(recipient.Channel(),fmt.Sprintf("CHAT: %d",chatRoomRef))
				sendLine(recipient.Channel(),fmt.Sprintf("CLIENT_NAME: %s",client.Name()))
				sendLine(recipient.Channel(),fmt.Sprintf("MESSAGE: %s",messageLines[0]))
				for i := 1; i < len(messageLines); i++ {
					sendLine(recipient.Channel(),messageLines[i])
				}
				sendLine(recipient.Channel(),"")
			}
		}

		rr.done <- STATUS_SUCCESS_CONTINUE
	}
}
