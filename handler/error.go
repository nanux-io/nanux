package handler

// ErrorHandler defines the type of function to handle errors return by action
type ErrorHandler func(error) []byte
