package protocol

import(
	"strconv"
	)

type Helo struct {
	ip string
	port int
	queue chan *Exchange
}

func MakeHelo(ip string, port int, threadCount int) *Helo{
	e := &Helo{
		ip,
		port,
		make(chan *Exchange),
	}
	for i := 0; i < threadCount; i++ {
		go e.runLoop()
	}
	return e
}

func (e *Helo) Identifier() string {
	return "HELO"
}

func (e *Helo) Handle(request <-chan byte, response chan<- byte) <-chan int {
	done := make(chan int)
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
		rr.done <- 1
	}
}
