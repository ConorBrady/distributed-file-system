package protocol

import(
	"regexp"
	"fmt"
	"log"
	"github.com/conorbrady/distributed-file-system/locate"
	)

type LocateFileProtocol struct {
	queue chan *Exchange
}

func MakeLocateFileProtocol(threadCount int) *LocateFileProtocol {
	p := &LocateFileProtocol{
		make(chan *Exchange, threadCount),
	}
	for i := 0; i < threadCount; i++ {
		go p.runLoop()
	}
	return p
}

func (p *LocateFileProtocol)Identifier() string {
	return "LOCATE"
}

func (p *LocateFileProtocol)Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
	done := make(chan StatusCode, 1)
	p.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (p *LocateFileProtocol)runLoop() {
	for {
		rr := <- p.queue

		// Line 1 "LOCATE:"
		r1, _ := regexp.Compile("\\A\\s*(\\S+)\\s*\\z")
		matches1 := r1.FindStringSubmatch(readLine(rr.request))
		if len(matches1) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		address := locate.LocateFile(matches1[1])

		sendLine(rr.response,fmt.Sprintf("ADDRESS: %s",address))

		log.Println("Responded with: "+fmt.Sprintf("ADDRESS: %s",address))
		rr.done <- STATUS_SUCCESS_DISCONNECT
	}
}
