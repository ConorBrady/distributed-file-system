package protocol

import (
	"regexp"
	"distributed-file-system/auth"
	"encoding/base64"
	)

type ServiceTicketProtocol struct {
	queue chan *Exchange
}

func MakeServiceTicketProtocol(threadCount int) *ServiceTicketProtocol{
	e := &ServiceTicketProtocol{
		make(chan *Exchange,threadCount),
	}
	for i := 0; i < threadCount; i++ {
		go e.runLoop()
	}
	return e
}

func (e *ServiceTicketProtocol) Identifier() string {
	return "SERVICE_TICKET"
}

func (e *ServiceTicketProtocol) Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
	done := make(chan StatusCode)
	e.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (e *ServiceTicketProtocol) runLoop() {
	for {
		rr := <- e.queue

		// "SERVICE_TICKET:"
		r1, _ := regexp.Compile("\\A\\s*(\\w+=*)\\s*\\z")
		matches1 := r1.FindStringSubmatch(readLine(rr.request))
		if len(matches1) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		sessionKey := auth.GenFromTicket(matches1[1])

		if sessionKey == nil {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		// "AUTHENTICATOR:"
		r2, _ := regexp.Compile("\\AAUTHENTICATOR:\\s*(\\w+=*)\\s*\\z")
		matches2 := r2.FindStringSubmatch(readLine(rr.request))
		if len(matches2) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		sendLine(rr.response,"ENCRYPTED_SESSION_KEY: " + base64.StdEncoding.EncodeToString(sessionKey.EncryptedKey()))
		sendLine(rr.response,"SERVICE_TICKET: "+base64.StdEncoding.EncodeToString(sessionKey.MarshalAndEncrypt()))

		rr.done <- STATUS_SUCCESS_CONTINUE
	}
}
