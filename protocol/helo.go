package protocol

import(
	"strconv"
	)

type Helo struct {
	ip string
	port int
	queue chan *Exchange
	uuid string
	mode string
}

func MakeHelo(ip string, port int, threadCount int, mode string, uuid string) *Helo{
	e := &Helo{
		ip,
		port,
		make(chan *Exchange),
		uuid,
		mode,
	}
	for i := 0; i < threadCount; i++ {
		go e.runLoop()
	}
	return e
}

func (e *Helo) Identifier() string {
	return "HELO"
}

func (e *Helo) Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
	done := make(chan StatusCode)
	e.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (e *Helo) runLoop() {
	for {
		rr := <- e.queue
		for _, b := range []byte("HELO:") {
			rr.response <- b
		}
		for msgByte := range rr.request {
			if(msgByte=='\n') {
				break
			}
			rr.response <- msgByte
		}


		for _, b := range []byte("\nIP: "+e.ip+"\nPort: "+strconv.Itoa(e.port)+"\nStudentID: 08506426\nUUID: "+e.uuid+"\nMODE: "+e.mode+"\n") {
			rr.response <- b
		}
		rr.done <- STATUS_SUCCESS_CONTINUE
	}
}
