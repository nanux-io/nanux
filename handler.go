package nanux

// Handler is the fundamental building block of servers and clients.
// It represents a single RPC method and its options. The context correspond to the context
// provided to the Nanux instance. It is not a context request scoped
// but an app scoped one
type Handler struct {
	Fn   HandlerFunc
	Opts []HandlerOpt
}

// HandlerFunc is the signature of the functions used by a Handler
type HandlerFunc func(ctx *interface{}, request Request) (response []byte, err error)

// HandlerOpt is option for action. It will be used by the transporter. The Name
// is the name of the option it is defined at the transporter lvl and the
// value must be anything. Check the doc of the options to know which value
// is accepted for each option.
type HandlerOpt struct {
	Name  HandlerOptName
	Value interface{}
}

// HandlerOptName is name of an option. It might be used by transporters.
type HandlerOptName string

// THandler is the handler used by the transporter. The difference with a
// simple Handler is that there is no context in the function because it is provided by the
// Nanux instance (check the `Handle` method of Nanux)
type THandler struct {
	Fn   func(request Request) (response []byte, err error)
	Opts []HandlerOpt
}

// ErrorHandler defines the type of function to handle errors return by action
type ErrorHandler func(error, Request) []byte

// Request contains the core data of the request (eg: the body for an http request)
// and some extra data linked to the request (eg: metadata, input type json, and output type)
type Request struct {
	Data []byte
	M    map[string]interface{}
}
