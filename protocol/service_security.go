package protocol

import (
	"log"
	"regexp"
	"distributed-file-system/auth/service"
	"distributed-file-system/auth/crypto"
	"encoding/base64"
	"fmt"
	)

type ServiceSecurityProtocol struct {
	queue chan *Exchange
	router *Router
}

func MakeServiceSecurityProtocol(threadCount int) *ServiceSecurityProtocol{
	e := &ServiceSecurityProtocol{
		make(chan *Exchange,threadCount),
		MakeRouter(),
	}

	for i := 0; i < threadCount; i++ {
		go e.runLoop()
	}
	return e
}

func (e *ServiceSecurityProtocol) AddProtocol(p Protocol) {
	e.router.AddProtocol(p)
}

func (e *ServiceSecurityProtocol) Identifier() string {
	return "SERVICE_TICKET"
}

func (e *ServiceSecurityProtocol) Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
	done := make(chan StatusCode)
	e.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (e *ServiceSecurityProtocol) runLoop() {
	for {
		rr := <- e.queue

		fmt.Println("Started service ticket")
		// "SERVICE_TICKET:"
		r1, _ := regexp.Compile("\\A\\s*(\\S+)\\s*\\z")
		line1 := readLine(rr.request)
		matches1 := r1.FindStringSubmatch(line1)
		if len(matches1) < 2 {
			log.Println("Recieved "+line1)
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		encryptedSessionKey, encSesErr := base64.StdEncoding.DecodeString(matches1[1])

		if encSesErr != nil {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		sessionKey := service.DecryptSessionKey(encryptedSessionKey)

		if sessionKey == nil {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		log.Println("Waiting on authenticator")
		// "AUTHENTICATOR:"
		r2, _ := regexp.Compile("\\AAUTHENTICATOR:\\s*(\\S+)\\s*\\z")
		matches2 := r2.FindStringSubmatch(readLine(rr.request))
		if len(matches2) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}
		log.Println("Authenticator recieved")

		encryptedAuthenticator, encAuthErr := base64.StdEncoding.DecodeString(matches2[1])

		if encAuthErr != nil {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		log.Println("Decrypting authenticator")
		authenticator := service.DecryptAuthenticator(encryptedAuthenticator, sessionKey.Key())
		log.Println("Authenticator decrypted")

		if authenticator == nil {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		if authenticator.Username() != sessionKey.Username() {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		sendLine(rr.response,"RESPONSE: " + base64.StdEncoding.EncodeToString(authenticator.MakeResponse(sessionKey.Key())))

		requestChan := make(chan byte)
		responseChan := make(chan byte)

		stop := false

		go func() {
			for !stop {
				line := readLine(responseChan)
				enc := crypto.EncryptString(line,sessionKey.Key())
				for _, responseByte := range []byte(base64.StdEncoding.EncodeToString(enc)) {
					rr.response <- responseByte
					fmt.Print(responseByte)
				}
				rr.response <- '\n'
				fmt.Print('\n')
			}
		}()

		go func() {

			for !stop {
				encryptedData, _ := base64.StdEncoding.DecodeString(readLine(rr.request))

				for _, b := range []byte(crypto.DecryptToString(encryptedData,sessionKey.Key())) {
					requestChan <- b
				}
			}
		}()

		status := STATUS_UNDEFINED

		for status != STATUS_SUCCESS_DISCONNECT {

			log.Println("Secure channel entered")

			buffer := make([]byte,0)

			for nb := <- requestChan; nb != '\n' && nb != ' ' && nb != ':' && nb != '\r'; nb = <- requestChan {
				buffer = append(buffer,nb)
			}

			ident := string(buffer)

			if ident != "" {
				status = <- e.router.Route(ident,requestChan,responseChan)
			}
		}

		stop = true

		rr.done <- STATUS_SUCCESS_CONTINUE
	}
}
