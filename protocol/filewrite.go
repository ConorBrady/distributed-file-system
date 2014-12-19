package protocol

import(
	"regexp"
	"os"
	"io"
	"fmt"
	"encoding/hex"
	"crypto/sha256"
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

			// Line 1 "WRITE_FILE:"
			r1, _ := regexp.Compile("\\A\\s*(\\S+)\\s*\\z")
			matches1 := r1.FindStringSubmatch(readLine(rr.request))
			if len(matches1) < 2 {
				respondError(ERROR_MALFORMED_REQUEST,rr.response)
				rr.done <- STATUS_ERROR
				continue
			}

			hash := sha256.New()
			hash.Write([]byte(matches1[1]))
			md := hash.Sum(nil)
			mdStr := hex.EncodeToString(md)

			tempFileName := "tmp/"+mdStr
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

			for i := 0; i <= (contentLength >> 7); i++ {
				file.Write(readByteCount(rr.request,(contentLength-(i<<7))%(1<<7)))
			}
			// if read(rr.request,':') != "CONTENT_BASE64" {
			// 	respondError(ERROR_MALFORMED_REQUEST,rr.response)
			// 	rr.done <- STATUS_ERROR
			// 	continue
			// }
			//
			// decoder := base64.NewDecoder(base64.StdEncoding,MakeChannelReader(rr.request))
			//
			// buffer := make([]byte,48)
			// n, errDec := decoder.Read(buffer)
			//
			// for errDec == nil {
			// 	_, err := file.Write(buffer[:n])
			// 	fmt.Println("Writing "+string(buffer[:n]))
			// 	if err != nil {
			// 		fmt.Println(err.Error())
			// 	}
			// 	fmt.Println(n)
			// 	if n < len(buffer) {
			// 		break
			// 	}
			// 	n, errDec = decoder.Read(buffer)
			// }

			file.Close()

			if errDec != nil && errDec != io.EOF {
				respondError(ERROR_MALFORMED_REQUEST,rr.response)
				os.Remove(tempFileName)
				rr.done <- STATUS_ERROR
				continue
			}

			os.Rename(tempFileName,"storage/"+mdStr)

			sendLine(rr.response,"SUCCESS")

			rr.done <- STATUS_SUCCESS_CONTINUE
		}
	}
