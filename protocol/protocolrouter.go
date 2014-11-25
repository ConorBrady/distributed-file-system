package protocol

import(
	"log"
	)

type Router struct {
	protocols map[string]Protocol
}

func MakeRouter() *Router {
	return &Router{
		make(map[string]Protocol),
	}
}

func (r *Router) Route(identifier string, request <-chan byte, response chan<- byte) <-chan StatusCode {

	protocol, ok := r.protocols[identifier]

	if ok {
		return protocol.Handle(request,response)
	} else {
		error(ERROR_UNKNOWN_PROTOCOL,response)
		channel := make(chan StatusCode,1)
		channel <- STATUS_ERROR
		return channel
	}
}

func (r *Router) AddProtocol(protocol Protocol) {

	if _, taken := r.protocols[protocol.Identifier()]; taken {
		log.Fatal("Tried to add protocol "+protocol.Identifier()+" more than once")
	}

	r.protocols[protocol.Identifier()] = protocol

	_ , ok := r.protocols[protocol.Identifier()]
	if !ok {
		log.Fatal("Failed to add protocol "+protocol.Identifier())
	}
}
