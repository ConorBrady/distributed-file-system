package main

import(
    "flag"
    "log"
    "os"

    "github.com/conorbrady/distributed-file-system/chat"
    "github.com/conorbrady/distributed-file-system/protocol"
    "github.com/conorbrady/distributed-file-system/tcpserver"
    "github.com/conorbrady/distributed-file-system/auth/authentication"
    "github.com/conorbrady/distributed-file-system/locate"

    "code.google.com/p/go-uuid/uuid"
    )

func main (){

    logFile, _ := os.Create("out.log")
    log.SetOutput(logFile)

    port := flag.Int("port",-1,"Port to listen on")
    threadCount := flag.Int("threadCount", 4, "Available thread count")
    mode := flag.String("mode","","Server mode, select between 'DS', 'AS' and 'FS'")
    keyPath := flag.String("key","","Private key AES256 key path")

    flag.Parse()

    if *port<0 {
        log.Fatal("Must pass port via -port x flag")
    }

    if *mode == "" {
        log.Fatal("Please select mode from 'DS', 'AS' and 'FS' with -mode")
    }

    if *keyPath == "" {
        log.Fatal("Please provide private key with -key")
    }

    wd, _ := os.Getwd()
    uuidFilePath := wd+"/uuid.txt"

    uuidFile, fileErr := os.Open(uuidFilePath)

    if fileErr != nil {

        uuidFile, _ = os.Create(uuidFilePath)
        uuidFile.Write([]byte(uuid.New()))
        uuidFile.Close()

        uuidFile, fileErr = os.Open(uuidFilePath)

        if fileErr != nil {
            log.Fatal("File not writing")
        }
    }

    uuid := make([]byte,36)
    uuidFile.Read(uuid)

    tcpServer := tcpserver.MakeTCPServer(os.Getenv("IP_ADDRESS"),*port,*threadCount)

    chat := chat.MakeChat()

    tcpServer.AddProtocol(protocol.MakeHelo(os.Getenv("IP_ADDRESS"),*port,4,*mode,string(uuid)))

    tcpServer.AddProtocol(protocol.MakeChatJoinProtocol(chat,4))
    tcpServer.AddProtocol(protocol.MakeChatLeaveProtocol(chat,4))
    tcpServer.AddProtocol(protocol.MakeChatMessageProtocol(chat,4))
    tcpServer.AddProtocol(protocol.MakeDisconnectProtocol(chat,1))

    switch *mode {
        case "AS":
            tcpServer.AddProtocol(protocol.MakeAuthenticationProtocol(4,*keyPath))

            go authentication.RunManagement()

        case "DS":
            secureProtocol := protocol.MakeServiceSecurityProtocol(4,*keyPath)
            secureProtocol.AddProtocol(protocol.MakeLocateFileProtocol(4))

            tcpServer.AddProtocol(secureProtocol)

            go locate.RunManagement()

        case "FS":
            secureProtocol := protocol.MakeServiceSecurityProtocol(4,*keyPath)
            secureProtocol.AddProtocol(protocol.MakeFileReadProtocol(4))
            secureProtocol.AddProtocol(protocol.MakeFileWriteProtocol(4))

            tcpServer.AddProtocol(secureProtocol)

        default:
            log.Fatal("Invalid server mode selected, select from 'DS', 'AS' and 'FS'")
    }

    tcpServer.BlockingRun()
}
