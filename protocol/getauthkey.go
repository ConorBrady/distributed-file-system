package protocol

import(
	"regexp"
	"os"
	"fmt"
	"code.google.com/p/go-sqlite/go1/sqlite3"
	)

	type LoginUserProtocol struct {
		queue chan *Exchange
	}

	func MakeLoginUserProtocol(threadCount int) *LoginUserProtocol{
		e := &LoginUserProtocol{
			make(chan *Exchange,threadCount),
		}
		for i := 0; i < threadCount; i++ {
			go e.runLoop()
		}
		return e
	}

	func (e *LoginUserProtocol) Identifier() string {
		return "LOGIN_USERNAME"
	}

	func (e *LoginUserProtocol) Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
		done := make(chan StatusCode)
		e.queue <- &Exchange{
			request,
			response,
			done,
		}
		return done
	}

	func (e *LoginUserProtocol) runLoop() {
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

			db, dbErr := sqlite3.Open(os.Getenv("GOPATH")+"/src/distributed-file-system/db/users.sqlite")
			if dbErr!=nil {
				fmt.Println(dbErr.Error())
			}

			args := sqlite3.NamedArgs{"$username": matches1[1]}
			statement, qErr := db.Query("select password from users where username=$username", args)

			if qErr != nil {
				fmt.Println(qErr.Error())
				respondError(ERROR_USER_NOT_FOUND,rr.response)
				rr.done <- STATUS_ERROR
				continue
			}

			var password string
			statement.Scan(&password)
			sendLine(rr.response,"PASSWORD: "+password)

			rr.done <- STATUS_SUCCESS_CONTINUE
		}
	}
