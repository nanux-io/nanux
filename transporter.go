package nanux

// Transporter defines common methods to interact with different transporters
// (nats, http, nanomsg...) and use them for listening incoming request
type Transporter interface {
	Run() error
	Close() error
	Handle(route string, tHandler THandler) error
	HandleError(ErrorHandler) error
}
