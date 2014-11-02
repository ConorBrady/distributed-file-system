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

func (r *Router) Route(identifier string, request <-chan byte, response chan<- byte) {

	protocol, ok := r.protocols[identifier]

	if ok {
		log.Println("Found "+identifier+" protocol")
		protocol.Handle(request,response)
	} else {
		log.Println("Protocol "+identifier+" not found")
		log.Println("Possible protocols are:")
		for k, _ := range r.protocols {
			log.Println(k)
		}
		close(response)
	}
}

func (r *Router) AddProtocol(protocol Protocol) {
	if _, taken := r.protocols[protocol.Identifier()]; taken {
		log.Fatal("Tried to add protocol "+protocol.Identifier()+" more than once")
	}
	log.Println("Adding protocol "+protocol.Identifier())
	r.protocols[protocol.Identifier()] = protocol
	p, ok := r.protocols[protocol.Identifier()]
	if ok {
		log.Println("Successfully added "+p.Identifier())
	} else {
		log.Fatal("Failed to add protocol "+protocol.Identifier())
	}
}
