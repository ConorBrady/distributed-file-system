package protocol

import (
	"regexp"
	"log"
	"strconv"

	"github.com/conorbrady/distributed-file-system/file"
	)

type FileWriteProtocol struct {
	queue chan *Exchange
}

func MakeFileWriteProtocol(threadCount int) *FileWriteProtocol {
	p := &FileWriteProtocol{
		make(chan *Exchange, threadCount),
	}
	for i := 0; i < threadCount; i++ {
		go p.runLoop()
	}
	return p
}

func (p *FileWriteProtocol)Identifier() string {
	return "WRITE_FILE"
}

func (p *FileWriteProtocol)Handle(request <-chan byte, response chan<- byte) <-chan StatusCode {
	done := make(chan StatusCode, 1)
	p.queue <- &Exchange{
		request,
		response,
		done,
	}
	return done
}

func (p *FileWriteProtocol)runLoop() {
	for {
		rr := <- p.queue

		line := readLine(rr.request)
		log.Println("WRITE_FILE:"+line)
		// Line 1 "WRITE_FILE:"
		r1, _ := regexp.Compile("\\A\\s*(\\S+)\\s*\\z")
		matches1 := r1.FindStringSubmatch(line)
		if len(matches1) < 2 {
			log.Printf("Got string post WRITE FILE %q",line)
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		// Line 2 "START:"
		line = readLine(rr.request)
		log.Println(line)
		r2, _ := regexp.Compile("\\A\\s*START\\s*:\\s*(\\d+)\\s*\\z")
		matches2 := r2.FindStringSubmatch(line)
		if len(matches2) < 2 {
			log.Printf("Got string to WRITE FILE %q",line)
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		start, _ := strconv.Atoi(matches2[1])

		// Line 3 "CONTENT_LENGTH:"

		line = readLine(rr.request)
		log.Println(line)
		r3, _ := regexp.Compile("\\A\\s*CONTENT_LENGTH\\s*:\\s*(\\d+)\\s*\\z")
		matches3 := r3.FindStringSubmatch(line)
		if len(matches3) < 2 {
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		contentLength, _ := strconv.Atoi(matches3[1])

		syncToToken(rr.request,"DATA:")

		data := readByteCount(rr.request,contentLength)
		log.Println("Data recieved")
		log.Println(string(data))

		if err := file.WriteData(matches1[1],start,data); err != nil {
			log.Println(err.Error())
			respondError(ERROR_MALFORMED_REQUEST,rr.response)
			rr.done <- STATUS_ERROR
			continue
		}

		log.Println("Wrote "+matches1[1])

		block, _ := file.GetBlock(matches1[1],start/file.BlockSize)
		sendLine(rr.response, "WROTE_BLOCK:"+strconv.Itoa(start/file.BlockSize))
		sendLine(rr.response, "HASH: "+ block.Hash())

		rr.done <- STATUS_SUCCESS_DISCONNECT
	}
}
