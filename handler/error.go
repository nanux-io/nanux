package handler

// ManageError defines the type of function to handle errors return by action
type ManageError func(error) []byte
