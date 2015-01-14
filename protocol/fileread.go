package protocol

import (
	"regexp"
	"log"
	"strconv"

	"github.com/conorbrady/distributed-file-system/file"
)

type FileReadProtocol struct {
	queue chan *Exchange
}

func MakeFileReadProtocol(threadCount int) *FileReadProtocol {
	p := &FileReadProtocol{
		make(chan *Exchange, threadCount),
	}
	for i := 0; i < threadCount; i++ {
		go p.runLoop()
	}
	return p
}

func (p *FileReadProtocol)Identifier() string {
	return "READ_FILE"
}

func (p *FileReadProtocol)Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
	done := make(chan StatusCode, 1)
	p.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (p *FileReadProtocol)runLoop() {
	for {
		rr := <- p.queue

		line := readLine(rr.request)
		log.Println("READ_FILE:"+line)

		// Line 1 "READ_FILE:"
		r1, _ := regexp.Compile("\\A\\s*(\\S+)\\s*\\z")
		matches1 := r1.FindStringSubmatch(line)
		if len(matches1) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		// Line 2 "BLOCK_INDEX:"
		r2, _ := regexp.Compile("\\AOFFSET:\\s*(\\d+)\\s*\\z")
		matches2 := r2.FindStringSubmatch(readLine(rr.request))
		if len(matches2) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		offset, _ := strconv.Atoi(matches2[1])
		blockIndex := offset/file.BlockSize

		data, err := file.ReadData(matches1[1],blockIndex)

		if err != nil {
			log.Println(err.Error())
			respondError(ERROR_FILE_NOT_FOUND,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		block, _ := file.GetBlock(matches1[1],blockIndex)

		sendLine(rr.response, "START: "+ strconv.Itoa(blockIndex*file.BlockSize))
		sendLine(rr.response, "HASH: "+ block.Hash())
		sendLine(rr.response, "CONTENT_LENGTH: "+ strconv.Itoa(len(data)))

		for _, b := range data {
			rr.response <- b
		}

		rr.done <- STATUS_SUCCESS_DISCONNECT
	}
}
