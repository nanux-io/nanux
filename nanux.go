package nanux

// Nanux is the base struct to use nanux. It contains the used listeners
// (nats, http, ...) and a context which can hold for example the instance
// of the DB to used.
type Nanux struct {
	// L is the instance of listener used by nanux
	T Transporter
	// Ctx is an a nanux scoped context
	Ctx          interface{}
	errorHandler ErrorHandler
}

// Handle defines the function to execute when the given route is reached on the
// listener.
// A route is an http route in the case of an http listener, its a channel
// subscription in the case of a nats listener etc...
func (n *Nanux) Handle(route string, handler Handler, middlewares ...Middleware) error {
	fn := func(req Request) ([]byte, error) {
		return chainMiddleware(handler.Fn, middlewares...)(&n.Ctx, req)
	}

	tHandler := THandler{
		Fn:   fn,
		Opts: handler.Opts,
	}
	err := n.T.Handle(route, tHandler)

	return err
}

// Run launch the transporter (nats, http, etc...) attached to
// Nanux.
func (n *Nanux) Run() error {
	return n.T.Run()
}

// Close the listener connection
func (n *Nanux) Close() error {
	return n.T.Close()
}

// HandleError specify error handler which must be called when a handler return
// en error
func (n *Nanux) HandleError(errHandler ErrorHandler) error {
	return n.T.HandleError(errHandler)
}

// New create a new Nanux for nanux
func New(transporter Transporter, ctx interface{}) *Nanux {
	return &Nanux{
		T:   transporter,
		Ctx: ctx,
	}
}
