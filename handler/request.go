package handler

// Request contains the core data of the request (eg: the body for an http request)
// and some metadata (eg: input type json, and output type)
type Request struct {
	Data []byte
	Meta map[string]interface{}
}
