package protocol

import (
	"regexp"
	"log"
	"distributed-file-system/auth/authentication"
	"encoding/base64"
	)

type AuthenticationProtocol struct {
	queue chan *Exchange
	keyPath string
}

func MakeAuthenticationProtocol(threadCount int, keyPath string) *AuthenticationProtocol{
	e := &AuthenticationProtocol{
		make(chan *Exchange,threadCount),
		keyPath,
	}
	for i := 0; i < threadCount; i++ {
		go e.runLoop()
	}
	return e
}

func (e *AuthenticationProtocol) Identifier() string {
	return "AUTHENTICATE"
}

func (e *AuthenticationProtocol) Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
	done := make(chan StatusCode)
	e.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (e *AuthenticationProtocol) runLoop() {
	for {
		rr := <- e.queue

		line := readLine(rr.request)
		log.Println("AUTHENTICATE:"+line)
		// "AUTHENTICATE:"
		r1, _ := regexp.Compile("\\A\\s*(\\w+)\\s*\\z")
		matches1 := r1.FindStringSubmatch(line)
		if len(matches1) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		user := authentication.GetUser(matches1[1])

		if user == nil {
			respondError(ERROR_USER_NOT_FOUND,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		sessionKey := authentication.GetSessionKey(*user)

		sendLine(rr.response,"ENCRYPTED_SESSION_KEY: " + base64.StdEncoding.EncodeToString(sessionKey.EncryptedKey()))
		sendLine(rr.response,"SERVICE_TICKET: "+base64.StdEncoding.EncodeToString(sessionKey.MarshalAndEncrypt(e.keyPath)))

		rr.done <- STATUS_SUCCESS_CONTINUE
	}
}
