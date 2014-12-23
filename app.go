package main

import(
    "flag"
    "log"
    "os"

    "distributed-file-system/chat"
    "distributed-file-system/protocol"
    "distributed-file-system/tcpserver"
    "distributed-file-system/auth/authentication"
    )

func main (){

    port := flag.Int("port",-1,"Port to listen on")
    threadCount := flag.Int("threadCount", 100, "Available thread count")
    mode := flag.String("mode","","Server mode, select between 'DS', 'AS' and 'FS'")

    flag.Parse()

    if *port<0 {
        log.Fatal("Must pass port via -port x flag")
    }

    if *mode == "" {
        log.Fatal("Please select mode from 'DS', 'AS' and 'FS'")
    }

    tcpServer := tcpserver.MakeTCPServer(os.Getenv("IP_ADDRESS"),*port,*threadCount)

    chat := chat.MakeChat()

    tcpServer.AddProtocol(protocol.MakeHelo(os.Getenv("IP_ADDRESS"),*port,4))

    tcpServer.AddProtocol(protocol.MakeChatJoinProtocol(chat,4))
    tcpServer.AddProtocol(protocol.MakeChatLeaveProtocol(chat,4))
    tcpServer.AddProtocol(protocol.MakeChatMessageProtocol(chat,4))
    tcpServer.AddProtocol(protocol.MakeDisconnectProtocol(chat,1))

    switch *mode {
        case "AS":
            tcpServer.AddProtocol(protocol.MakeAuthenticationProtocol(4))

            go authentication.RunManagement()

        case "DS":

        case "FS":
            secureProtocol := protocol.MakeServiceSecurityProtocol(4)
            secureProtocol.AddProtocol(protocol.MakeFileReadProtocol(4))
            secureProtocol.AddProtocol(protocol.MakeFileWriteProtocol(4))

            tcpServer.AddProtocol(secureProtocol)

        default:
            log.Fatal("Incorrect server mode selected, select from 'DS', 'AS' and 'FS'")
    }

    tcpServer.BlockingRun()





}
