package gametime

type Logger interface {
	Info(...interface{})
	Error(...interface{})
}
