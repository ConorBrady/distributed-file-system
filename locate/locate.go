package locate

import (
	"math/rand"
	)

func LocateFile(filename string) string {

	file := getFile(filename)

	if file != nil {

		return getFileServer(file.fileServerUUID).address

	} else {

		fss := getFileServers()
		fs := fss[rand.Int()%len(fss)]

		file := createFile(filename,fs)

		return getFileServer(file.fileServerUUID).address
	}
}
