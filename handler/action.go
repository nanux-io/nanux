// Package handler define the type action and its parameters
package handler

// Action is the fundamental building block of servers and clients.
// It represents a single RPC method and its options. The context correspond to the context
// provided to the Nanux instance. It is not a context request scoped
// but an app scoped one
type Action struct {
	Fn   func(ctx *interface{}, request Request) (response []byte, err error)
	Opts []Opt
}

// Opt is option for action. It will be used by the transporter. The Name
// is the name of the option it is defined at the transporter lvl and the
// value must be anything. Check the doc of the options to know which value
// is accepted for each option.
type Opt struct {
	Name  OptName
	Value interface{}
}

// OptName is name of an option. It might be used by transporters.
type OptName string

// ActionListener is the action used by the transporter. The difference with a
// simple Action is that there is no context in the function because it is provided by the
// Nanux instance (check the `On` method of Nanux)
type ActionListener struct {
	Fn   func(request Request) (response []byte, err error)
	Opts []Opt
}
