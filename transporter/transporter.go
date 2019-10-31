package transporter

import (
	"github.com/nanux-io/nanux/handler"
)

// Listener defines common methods to interact with different transporters
// (nats, http, nanomsg...) and use them for listening incoming request
type Listener interface {
	Listen() error
	Close() error
	HandleAction(route string, action handler.ActionListener) error
	HandleError(handler.ErrorHandler) error
}
