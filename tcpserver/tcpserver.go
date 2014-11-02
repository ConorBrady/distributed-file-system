package tcpserver


import(
    "net"
    "bufio"
    "log"
    "strconv"
    "../protocol"
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
                server.sharedChan <- tcpConn
            case <- server.killChan:
                return
        }
    }
}

func (server *TCPServer) connectionHandler() {

    for {

        tcpConn := <- server.sharedChan
        reader := bufio.NewReader(tcpConn)

        buffer := make([]byte,0)

        for nb, _ := reader.ReadByte(); nb != '\n' && nb != ' ' && nb != '\r'; {
            buffer = append(buffer,nb)
            nb, _ = reader.ReadByte()
        }

        ident := string(buffer)

        requestChan := make(chan byte)
        responseChan := make(chan byte)
        if ident == "KILL_SERVICE" {
            server.killChan <- 1
            log.Println("Killing service")
            return
        }
        log.Println("Routing with ident "+strconv.Quote(ident))
        server.router.Route(ident,requestChan,responseChan)

        go func(){
            for nb, err := reader.ReadByte(); err==nil; nb, err = reader.ReadByte(){
                requestChan <- nb
            }
            log.Println("Finished receiving")
        }()
        log.Println("Responding")
        for responseByte := range responseChan {
            tcpConn.Write([]byte{responseByte})
        }
        log.Println("Finished Responding")

        tcpConn.Close()
    }
}
