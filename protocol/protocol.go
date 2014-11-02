package protocol


type Protocol interface {
	Identifier() string
	Handle(request <-chan byte, response chan<- byte)
}

type ProtocolPair struct {
	request <-chan byte
	response chan<- byte
}
