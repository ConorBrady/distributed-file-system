package locate

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	)
type File struct {
	filename string
	fileServerUUID string
}

func getFile(name string) *File{

	args := sqlite3.NamedArgs{
		"$filename": 			name,
	}

	query, qErr := dbConnect().Query("select file_server_uuid from files where name=$filename", args)

	if qErr != nil {
		return nil
	}

	var fs_uuid string

	query.Scan(&fs_uuid)
	query.Close()

	return &File{
		name,
		fs_uuid,
	}
}

func getFilesForFileServer(fs FileServer) []File {

	f := make([]File, 0)

	args := sqlite3.NamedArgs{
		"$file_server_uuid": 	fs.uuid,
	}

	query, qErr := dbConnect().Query("select name from files where file_server_uuid=$file_server_uuid",args)

	for qErr == nil {

		var name string

		query.Scan(&name)

		f = append(f, File{
			name,
			fs.uuid,
		})

		qErr = query.Next()
	}

	return f
}

func createFile(name string, fs FileServer) *File{

	args := sqlite3.NamedArgs{
		"$file_server_uuid": 	fs.uuid,
		"$name":				name,
	}

	if err := dbConnect().Exec("insert into files ( file_server_uuid, name ) values ( $file_server_uuid, $name )", args); err != nil {
		return nil
	}

	return getFile(name)
}
