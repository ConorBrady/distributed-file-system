package tcpserver


import(
    "net"
    "strings"
    "bufio"
    "log"
    "strconv"
    )

type Protocol struct {
    Initiator func(request string) (shouldRun bool)
    Handler func(request string) (response string, kill bool)
}

type TCPServer struct {
    protocols []Protocol
    tcpAddr *net.TCPAddr
    threadCount int
}

func New(ip string, port int, threadCount int) *TCPServer{
    tcpAddr, _ := net.ResolveTCPAddr("tcp", ip+":"+strconv.Itoa(port))
    return &TCPServer{
        protocols:make([]Protocol,0),
        tcpAddr:tcpAddr,
        threadCount:threadCount,
    }
}

func (server* TCPServer) AddProtocol(protocol Protocol) {
    server.protocols = append(server.protocols,protocol)
}

func (server* TCPServer) BlockingRun() {

    tcpListener, err := net.ListenTCP("tcp", server.tcpAddr)

    if err != nil {
        log.Fatal(err)
    }

    sharedChan := make(chan *net.TCPConn, server.threadCount)
    killChan := make(chan int)
    tcpChan := make(chan *net.TCPConn)

    go func(){
        for {
            tcpConn, _ := tcpListener.AcceptTCP()
            tcpChan <- tcpConn
        }
    }()

    for i := 0; i < server.threadCount; i++ {
        go connectionHandler(sharedChan,killChan,server.protocols)
    }

    for {
        select {
            case tcpConn := <- tcpChan:
                select {
                    case sharedChan <- tcpConn:
                    default:
                }
            case <- killChan:
                return
        }
    }
}

func connectionHandler(sharedChan chan *net.TCPConn, killChan chan int, protocols []Protocol) {

    for {

        tcpConn := <- sharedChan
        message, _ := bufio.NewReader(tcpConn).ReadString('\n')
        message = strings.TrimSpace(message)

        for _, protocol := range protocols {
            if protocol.Initiator(message) {
                response, kill := protocol.Handler(message)
                tcpConn.Write([]byte(response))
                if kill {
                    killChan <- 1
                }
            }
        }

        tcpConn.Close()
    }
}
