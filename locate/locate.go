package locate

import (
	"math/rand"
	"log"
	)

func LocateFile(filename string) string {

	log.Println("Locating "+filename)

	if file := getFile(filename); file != nil {

		fs := getFileServer(file.fileServerUUID)

		if *FSConnect(fs.address).getUUID() == file.fileServerUUID { // Is that server at that address

			log.Println("File found at "+fs.address)

			return fs.address

		} else {

			log.Println("Server gone, deleting records")

			deleteFileServer(fs.uuid)
		}
	}

	log.Println("Allocating file location")

	fss := getFileServers()
	fs := fss[rand.Int()%len(fss)]

	for FSConnect(fs.address) == nil {

		deleteFileServer(fs.uuid)
		fss = getFileServers()
		fs = fss[rand.Int()%len(fss)]
	}

	file := createFile(filename,fs)

	return getFileServer(file.fileServerUUID).address
}
