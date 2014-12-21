package protocol

import(
	"regexp"
	"distributed-file-system/auth"
	)

	type AuthenticationProtocol struct {
		queue chan *Exchange
	}

	func MakeAuthenticationProtocol(threadCount int) *AuthenticationProtocol{
		e := &AuthenticationProtocol{
			make(chan *Exchange,threadCount),
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

			// Line 1 "LOGIN_USERNAME:"
			r1, _ := regexp.Compile("\\A\\s*(\\w+)\\s*\\z")
			matches1 := r1.FindStringSubmatch(readLine(rr.request))
			if len(matches1) < 2 {
				respondError(ERROR_MALFORMED_REQUEST,rr.response)
				rr.done <- STATUS_ERROR
				continue
			}

			user := auth.GetUser(matches1[1])

			if user == nil {
				respondError(ERROR_USER_NOT_FOUND,rr.response)
				rr.done <- STATUS_ERROR
				continue
			}

			sessionKey := auth.GetSessionKey(*user)
			sendLine(rr.response,"KEY: "+sessionKey.Key())

			rr.done <- STATUS_SUCCESS_CONTINUE
		}
	}
