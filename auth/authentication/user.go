package authentication

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"crypto/sha256"
	"strings"
	)

type User struct {
	userId int
	username string
	password string
}

func GetUser(username string) *User {
	return getUser(username)
}

func getUsers() []User {

	users := make([]User, 0)
	query, qErr := dbConnect().Query("select user_id, username, password from users")

	for qErr == nil {

		var password string
		var username string
		var userId int

		query.Scan(&userId,&username,&password)

		users = append(users, User{
			userId,
			username,
			password,
		})

		qErr = query.Next()
	}

	return users
}

func getUser(username string) *User{

	args := sqlite3.NamedArgs{
		"$username": username,
	}

	query, qErr := dbConnect().Query("select user_id, password from users where username=$username", args)

	if qErr != nil {
		return nil
	}

	var password string
	var userId int

	query.Scan(&userId,&password)
	query.Close()

	return &User{
		userId,
		username,
		password,
	}
}

func createUser(username string, password string) *User {

	if username != "" && password != "" {

		args := sqlite3.NamedArgs{
			"$username": username,
			"$password": password,
		}

		dbConnect().Exec("insert into users ( username, password ) values ( $username, $password )", args)

		return getUser(username)

	} else {

		return nil
	}

}

func deleteUser(username string) {

	user := getUser(username)

	if user != nil {
		args := sqlite3.NamedArgs{
			"$user_id": user.userId,
		}

		dbConnect().Exec("delete from users where user_id = $user_id", args)

		dbConnect().Exec("delete from session_keys where user_id = $user_id",args)
	}
}

func (u* User)getAESKey() []byte {
	key := sha256.Sum256([]byte(strings.TrimSpace(u.password)))
	return key[:]
}
