package nanux

// Middleware is the signature which must be fullfill by middlewares
type Middleware func(HandlerFunc) HandlerFunc

func chainMiddleware(f HandlerFunc, middlewares ...Middleware) HandlerFunc {
	// Check if there is some moddlewares
	if len(middlewares) == 0 {
		return f
	}

	// Pop out last middleware
	lenMws := len(middlewares)
	lastMiddleware := middlewares[lenMws-1]
	newMiddlewareList := middlewares[:lenMws-1]

	// Call chainMiddleware recursively
	return chainMiddleware(lastMiddleware(f), newMiddlewareList...)
}
