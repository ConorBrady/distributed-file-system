package protocol

import(
	"strconv"
	"log"
	)

type Helo struct {
	ip string
	port int
	queue chan *ProtocolPair
}

func MakeHelo(ip string, port int) *Helo{
	e := &Helo{
		ip,
		port,
		make(chan *ProtocolPair),
	}
	go e.runLoop()
	log.Println("HELO run loop running")
	return e
}

func (e *Helo) Identifier() string {
	return "HELO"
}

func (e *Helo) Handle(request <-chan byte, response chan<- byte) {
	e.queue <- &ProtocolPair{
		request,
		response,
	}
}

func (e *Helo) runLoop() {
	for {
		rr := <- e.queue
		for _, b := range []byte("HELO ") {
			rr.response <- b
		}
		for msgByte := range rr.request {
			if(msgByte=='\n') {
				break
			}
			rr.response <- msgByte
		}
		for _, b := range []byte("\nIP:"+e.ip+"\nPort:"+strconv.Itoa(e.port)+"\nStudentID:08506426\n") {
			rr.response <- b
		}
		close(rr.response)
	}
}
