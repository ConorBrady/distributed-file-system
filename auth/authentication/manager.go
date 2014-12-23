package authentication

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"text/tabwriter"
	"strconv"
	)

func RunManagement() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Please Select:")
		fmt.Println("1. View Users")
		fmt.Println("2. Add User")
		fmt.Println("3. Delete User")
		sel, _ := reader.ReadString('\n')
		sel = strings.TrimSpace(sel)
		switch sel {
			case "1":

				fmt.Println("*** Users on System ***")
				w := new(tabwriter.Writer)

				w.Init(os.Stdout, 0, 8, 0, '\t', 0)
				fmt.Fprintln(w,"ID\tUSERNAME")
				fmt.Println("----------------------------------")
				for _, u := range getUsers() {
					fmt.Fprintln(w,strconv.Itoa(u.userId)+"\t"+u.username)
				}
				w.Flush()

			case "2":
				fmt.Print("Username: ")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)

				fmt.Print("Password: ")
				password, _ := reader.ReadString('\n')
				password = strings.TrimSpace(password)

				if createUser(username,password) != nil {
					fmt.Println("User Added")
				} else {
					fmt.Println("An error occured")
				}
			case "3":
				fmt.Print("Username: ")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)

				deleteUser(username)
			default:
				fmt.Println("Please make a valid numerical choice")
		}

		fmt.Println()
	}

}
