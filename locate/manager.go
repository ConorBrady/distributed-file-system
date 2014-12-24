package locate

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"strings"
	"text/tabwriter"
	)

func RunManagement() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Please Select:")
		fmt.Println("1. View File Servers")
		fmt.Println("2. Add File Server")

		sel, _ := reader.ReadString('\n')
		sel = strings.TrimSpace(sel)
		switch sel {
			case "1":

				fmt.Println("*** File Servers on System ***")
				w := new(tabwriter.Writer)

				w.Init(os.Stdout, 0, 8, 0, '\t', 0)
				fmt.Fprintln(w,"UUID\tADDRESS")
				fmt.Println("----------------------------------")
				for _, fs := range getFileServers() {
					fmt.Fprintln(w,fs.uuid+"\t"+fs.address)
				}
				w.Flush()

			case "2":

				fmt.Print("Address: ")

				address, _ := reader.ReadString('\n')
				address = strings.TrimSpace(address)

				fmt.Println()

				log.Println("Connecting to server")

				if sc := FSConnect(address); sc != nil {

					log.Println("Server reached")

					uuid := sc.getUUID()

					if fs := getFileServer(*uuid); fs != nil {

						log.Println("Server already in directory")

						if err := fs.setAddress(address); err == nil {

							log.Println("Address updated")

						} else {

							log.Println(err.Error())
						}
					} else {

						log.Println("Adding server to directory")

						if createFileServer(address,*uuid) {

							log.Println("Fileserver Added")

						} else {

							log.Println("An error occured adding the file server")
						}
					}
				} else {

					log.Println("Server cannot be reached")
				}

			default:
				fmt.Println("Please make a valid numerical choice")
		}

		fmt.Println()
	}

}
