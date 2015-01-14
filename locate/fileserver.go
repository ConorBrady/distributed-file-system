package locate

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	)

type FileServer struct {
	address string
	uuid string
}

func getFileServers() []FileServer{

	fs := make([]FileServer, 0)

	query, qErr := dbConnect().Query("select address, uuid from file_servers")

	for qErr == nil {

		var address string
		var uuid string

		query.Scan(&address,&uuid)

		fs = append(fs, FileServer{
			address,
			uuid,
		})

		qErr = query.Next()
	}

	return fs
}

func getFileServer(uuid string) *FileServer {

	args := sqlite3.NamedArgs{
		"$uuid": uuid,
	}

	query, qErr := dbConnect().Query("select address from file_servers where uuid=$uuid", args)

	if qErr != nil {
		return nil
	}

	var address string

	query.Scan(&address)
	query.Close()

	return &FileServer{
		address,
		uuid,
	}
}

func (fs* FileServer)setAddress(address string) error {

	args := sqlite3.NamedArgs{
		"$address":	address,
		"$uuid":	fs.uuid,
	}

	err := dbConnect().Exec("update file_servers set address=$address where uuid=$uuid", args)

	if err == nil {
		fs.address = address
	}

	return err
}

func createFileServer(address string, uuid string) bool {

	args := sqlite3.NamedArgs{
		"$address":	address,
		"$uuid":	uuid,
	}

	err := dbConnect().Exec("insert into file_servers ( address, uuid ) values ( $address, $uuid )", args)

	return err == nil
}


func deleteFileServer(uuid string) {

	args := sqlite3.NamedArgs{
		"$uuid": uuid,
	}

	dbConnect().Exec("delete from file_servers where uuid = $uuid", args)

}
