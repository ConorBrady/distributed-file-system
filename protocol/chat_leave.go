package protocol

import(
	"regexp"
	"distributed-file-system/chat"
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
		make(chan *Exchange, threadCount),
	}
	for i := 0; i < threadCount; i++ {
		go p.runLoop()
	}
	return p
}

func (p *ChatLeaveProtocol)Identifier() string {
	return "LEAVE_CHATROOM"
}

func (p *ChatLeaveProtocol)Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
	done := make(chan StatusCode, 1)
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

		client.Invalidate()

		sendLine(rr.response,fmt.Sprintf("LEFT_CHATROOM: %d",chatRoomRef))
		sendLine(rr.response,fmt.Sprintf("JOIN_ID: %d",joinId))

		rr.done <- STATUS_SUCCESS_CONTINUE
	}
}
