package protocol


type Protocol interface {
	Identifier() string
	Handle(request <-chan byte, response chan<- byte) <-chan int
}

type Exchange struct {
	request <-chan byte
	response chan<- byte
	done chan<- int
}
