package gametime

type ErrorType uint8

const (
	_ ErrorType = iota
	Unrecoverable
	Recoverable
	Retriable
	Unknown
)

type Error struct {
	Actual error
	Type   ErrorType
}

func (e Error) Error() string {
	return e.Actual.Error()
}
