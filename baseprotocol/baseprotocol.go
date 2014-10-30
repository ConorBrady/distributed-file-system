package baseprotocol

import(
	"strconv"
	)

/**********************************/

type Echo struct {
	ip string
	port int
}

func MakeEcho(ip string, port int) *Echo{
	return &Echo{
		ip,
		port,
	}
}

func (e *Echo) Valid(request string) bool {
	return len(request) > 5 && request[:5] == "HELO "
}

func (e *Echo) Handle(request string) (response string, kill bool) {
	return request+"\nIP:"+e.ip+"\nPort:"+strconv.Itoa(e.port)+"\nStudentID:08506426\n", false
}


/***********************************/


type Kill struct {

}

func MakeKill() *Kill {
	return &Kill{}
}

func (k *Kill) Valid(message string) bool {
	return message == "KILL_SERVICE"
}

func (k *Kill) Handle(request string) (response string, kill bool) {
	return "", true
}
