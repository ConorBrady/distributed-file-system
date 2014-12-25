package protocol

import(
	"regexp"
	"os"
	"log"
	"net/url"
	"strconv"
	"code.google.com/p/go-uuid/uuid"
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
				respondError(ERROR_MALFORMED_REQUEST,rr.response)
				rr.done <- STATUS_ERROR
				continue
			}

			if err := os.Mkdir("tmp", 0777); err != nil {
				log.Println(err.Error())
			}

			if err := os.Mkdir("storage", 0777); err != nil {
				log.Println(err.Error())
			}

			tempFileName := "tmp/"+uuid.New()
			log.Println(tempFileName)
			file, _ := os.Create(tempFileName)

			// Line 2 "CONTENT_LENGTH:"
			r2, _ := regexp.Compile("\\ACONTENT_LENGTH:\\s*(\\d+)\\s*\\z")
			matches2 := r2.FindStringSubmatch(readLine(rr.request))
			if len(matches2) < 2 {
				respondError(ERROR_MALFORMED_REQUEST,rr.response)
				rr.done <- STATUS_ERROR
				continue
			}
			contentLength, _ := strconv.Atoi(matches2[1])

			// Body, read contentLength bytes

			for i := 0; i < (contentLength / 128); i++ {
				file.Write(readByteCount(rr.request,128))
			}

			file.Write(readByteCount(rr.request,contentLength%128))

			file.Close()

			os.Rename(tempFileName,"storage/"+url.QueryEscape(matches1[1]))

			log.Println("Wrote "+matches1[1])

			sendLine(rr.response,"SUCCESS")

			rr.done <- STATUS_SUCCESS_CONTINUE
		}
	}
