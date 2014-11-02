package main

import(
    "flag"
    "log"
    "os"
    "./tcpserver"
    "./protocol"
    )

func main (){

    port := flag.Int("port",-1,"Port to listen on")
    threadCount := flag.Int("threadCount", 100, "Available thread count")

    flag.Parse()

    if *port<0 {
        log.Fatal("Must pass port via -port x flag")
    }

    tcpServer := tcpserver.MakeTCPServer(os.Getenv("IP_ADDRESS"),*port,*threadCount)

    tcpServer.AddProtocol(protocol.MakeHelo(os.Getenv("IP_ADDRESS"),*port))

    tcpServer.BlockingRun()
}
