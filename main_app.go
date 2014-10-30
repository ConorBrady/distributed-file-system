package main


import(
    "flag"
    "log"
    "os"
    "strconv"
    "./tcpserver"
    )

func main (){

    port := flag.Int("port",-1,"Port to listen on")
    threadCount := flag.Int("threadCount", 100, "Available thread count")

    flag.Parse()

    if *port<0 {
        log.Fatal("Must pass port via -port x flag")
    }

    tcpServer := tcpserver.New(os.Getenv("IP_ADDRESS"),*port,*threadCount)


    tcpServer.AddProtocol(tcpserver.Protocol{
        Initiator:func(message string) bool {
            return len(message) > 5 && message[:5] == "HELO "
        },
        Handler:func(message string) (string, bool) {
            return message+"\nIP:"+os.Getenv("IP_ADDRESS")+"\nPort:"+strconv.Itoa(*port)+"\nStudentID:08506426\n", false
        },
    })

    tcpServer.AddProtocol(tcpserver.Protocol{
        Initiator:func(message string) bool {
            return message == "KILL_SERVICE"
        },
        Handler:func(message string) (string, bool) {
            return "", true
        },
    })

    tcpServer.BlockingRun()
}
