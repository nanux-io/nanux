package nanux

// Handler is the fundamental building block of servers and clients.
// It represents a single RPC method and its options. The context correspond to the context
// provided to the Nanux instance. It is not a context request scoped
// but an app scoped one
type Handler struct {
	Fn   HandlerFunc
	Opts map[HandlerOptName]interface{}
}

// HandlerFunc is the signature of the functions used by a Handler
type HandlerFunc func(ctx *interface{}, request Request) (response []byte, err error)

// HandlerOptName is name of an option. It might be used by transporters.
type HandlerOptName string

// THandler is the handler used by the transporter. The difference with a
// simple Handler is that there is no context in the function because it is provided by the
// Nanux instance (check the `Handle` method of Nanux)
type THandler struct {
	Fn   func(request Request) (response []byte, err error)
	Opts map[HandlerOptName]interface{}
}

// ErrorHandler defines the type of function to handle errors return by action
type ErrorHandler func(error, Request) []byte

// Request contains the core data of the request (eg: the body for an http request)
// and some extra data linked to the request (eg: metadata, input type json, and output type)
type Request struct {
	Data []byte
	M    map[string]interface{}
}
