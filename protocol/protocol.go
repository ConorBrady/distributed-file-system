package protocol

type StatusCode int

const (
	STATUS_SUCCESS_CONTINUE StatusCode = 0
	STATUS_SUCCESS_DISCONNECT StatusCode = 1
	STATUS_UNDEFINED StatusCode = 3
	STATUS_ERROR StatusCode = 4
)

type ErrorCode int

const (
	ERROR_MALFORMED_REQUEST ErrorCode = 0
	ERROR_UNKNOWN_PROTOCOL ErrorCode = 1
	ERROR_ROOM_NOT_FOUND ErrorCode = 2
	ERROR_CLIENT_NAME_MISMATCH ErrorCode = 3
	ERROR_JOIN_ID_NOT_FOUND ErrorCode = 4
	ERROR_FILE_NOT_FOUND = 5
	ERROR_USER_NOT_FOUND = 6
)

type Protocol interface {
	Identifier() string
	Handle(request <-chan byte, response chan<- byte) <-chan StatusCode
}

type Exchange struct {
	request <-chan byte
	response chan<- byte
	done chan<- StatusCode
}
