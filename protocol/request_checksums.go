package protocol

import (
	"regexp"
	"strconv"
	"github.com/conorbrady/distributed-file-system/file"
	)

type RequestChecksumsProtocol struct {
	queue chan *Exchange
}

func MakeRequestChecksumsProtocol(threadCount int) *RequestChecksumsProtocol {
	p := &RequestChecksumsProtocol{
		make(chan *Exchange, threadCount),
	}
	for i := 0; i < threadCount; i++ {
		go p.runLoop()
	}
	return p
}

func (p *RequestChecksumsProtocol)Identifier() string {
	return "REQUEST_CHECKSUMS"
}

func (p *RequestChecksumsProtocol)Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
	done := make(chan StatusCode, 1)
	p.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (p *RequestChecksumsProtocol)runLoop() {
	for {

		rr := <- p.queue

		// Line 1 "LOCATE:"
		r, _ := regexp.Compile("\\A\\s*(\\S+)\\s*\\z")
		matches := r.FindStringSubmatch(readLine(rr.request))
		if len(matches) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		filename := matches[1]

		sendLine(rr.response,"CHECKSUMS: "+filename)

		for {
			r, _ := regexp.Compile("\\AINDEX\\s*:\\s*(\\d+)\\s*\\z")
			matches := r.FindStringSubmatch(readLine(rr.request))
			if len(matches) < 2 {
				break
			}

			index, _ := strconv.Atoi(matches[1])

			block, _ := file.GetBlock(filename,index)
			sendLine(rr.response, "INDEX: "+strconv.Itoa(index))
			sendLine(rr.response, "HASH: "+block.Hash())
		}

		sendLine(rr.response,"END_CHECKSUMS")

		rr.done <- STATUS_SUCCESS_DISCONNECT
	}
}
