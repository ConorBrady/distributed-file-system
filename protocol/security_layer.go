package protocol

import (
	"log"
	"regexp"
	"distributed-file-system/auth/service"
	"distributed-file-system/auth/crypto"
	"encoding/base64"
)

type ServiceSecurityProtocol struct {
	queue chan *Exchange
	router *Router
	keyPath string
}

func MakeServiceSecurityProtocol(threadCount int, keyPath string) *ServiceSecurityProtocol{
	e := &ServiceSecurityProtocol{
		make(chan *Exchange,threadCount),
		MakeRouter(),
		keyPath,
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

		log.Println("Started service ticket")

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

		sessionKey := service.DecryptSessionKey(encryptedSessionKey,e.keyPath)

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

		log.Println("Connection secured")

		// HANDLES ENCRYPTION HERE

		requestChan := make(chan byte,32)
		responseChan := make(chan byte,32)

		stop := false

		go func() {
			for !stop {
				data := make([]byte,16)
				for i, _ := range data {
					b := <-responseChan
					data[i] = b

					if b == '\n' {
						break
					}
				}
				enc := crypto.EncryptBytes(data,sessionKey.Key())
				sendLine(rr.response,base64.StdEncoding.EncodeToString(enc)+"\n")
			}
		}()

		go func() {

			for !stop {
				enc, _ := base64.StdEncoding.DecodeString(readLine(rr.request))
				for _, b := range crypto.DecryptToBytes(enc,sessionKey.Key()) {

					requestChan <- b
				}
			}
		}()

		status := STATUS_UNDEFINED

		for status != STATUS_SUCCESS_DISCONNECT {

			log.Println("Parsing secure channel")

			buffer := make([]byte,0)

			log.Println("Initialized buffer")

			for nb := <- requestChan; nb != '\n' && nb != ' ' && nb != ':' && nb != '\r'; nb = <- requestChan {
				log.Println("Found "+string(nb))
				buffer = append(buffer,nb)
			}

			log.Println("Filled buffer")

			ident := string(buffer)

			log.Println("Passing "+ident+" to protocol router")

			if ident != "" {
				status = <- e.router.Route(ident,requestChan,responseChan)
			} else {
				log.Println("\"\" ident found")
			}
		}

		stop = true

		rr.done <- STATUS_SUCCESS_CONTINUE
	}
}
