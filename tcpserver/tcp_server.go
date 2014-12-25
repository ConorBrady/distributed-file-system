package tcpserver


import(
    "net"
    "bufio"
    "log"
    "strconv"
    "distributed-file-system/protocol"
    )

type TCPServer struct {
    router *protocol.Router
    tcpAddr *net.TCPAddr
    threadCount int
    killChan chan int
    sharedChan chan *net.TCPConn
}

func MakeTCPServer(ip string, port int, threadCount int) *TCPServer{
    tcpAddr, _ := net.ResolveTCPAddr("tcp", ip+":"+strconv.Itoa(port))
    return &TCPServer{
        router:protocol.MakeRouter(),
        tcpAddr:tcpAddr,
        threadCount:threadCount,
        killChan:make(chan int),
        sharedChan:make(chan *net.TCPConn, threadCount),
    }
}

func (server* TCPServer) AddProtocol(p protocol.Protocol) {
    server.router.AddProtocol(p)
}

func (server* TCPServer) BlockingRun() {

    tcpListener, err := net.ListenTCP("tcp", server.tcpAddr)

    if err != nil {
        log.Fatal(err)
    }

    tcpChan := make(chan *net.TCPConn)

    go func(){
        for {
            tcpConn, _ := tcpListener.AcceptTCP()
            tcpChan <- tcpConn
        }
    }()

    for i := 0; i < server.threadCount; i++ {
        go server.connectionHandler()
    }

    log.Println("Accepting connections")
    for {
        select {
            case tcpConn := <- tcpChan:
                select {
                    case server.sharedChan <- tcpConn:
                    default: // This is will drop any incoming connections if sharedChan is full
                }

            case <- server.killChan:
                return
        }
    }
}

func (server *TCPServer) connectionHandler() {

    for {

        tcpConn := <- server.sharedChan

        log.Println("New connection")

        reader := bufio.NewReader(tcpConn)
        responseChan := make(chan byte)
        requestChan := make(chan byte)

        go func(){
            for responseByte := range responseChan {
                tcpConn.Write([]byte{responseByte})
            }
        }()

        go func(){
            for nb, err := reader.ReadByte(); err==nil; nb, err = reader.ReadByte(){
                requestChan <- nb
            }
        }()

        status := protocol.STATUS_UNDEFINED

        for status != protocol.STATUS_SUCCESS_DISCONNECT {

            log.Println("Parsing for ident")

            buffer := make([]byte,0)

            for nb := <- requestChan; nb != '\n' && nb != ' ' && nb != ':' && nb != '\r'; nb = <- requestChan {
                buffer = append(buffer,nb)
            }

            ident := string(buffer)

            if ident != "" {
                if ident == "KILL_SERVICE" {
                    server.killChan <- 1
                    log.Println("Killing service")
                    return
                }
                status = <- server.router.Route(ident,requestChan,responseChan)
            } else {
                log.Println("\"\" ident found")
            }
        }

        tcpConn.Close()
    }
}
