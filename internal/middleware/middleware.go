package middleware

import "net/http"

// Middleware represents a function that wraps an http.Handler with additional
// functionality, such as logging, authentication, or request preprocessing.
type Middleware func(http.Handler) http.Handler

// WrapFunc wraps an http.HandlerFunc with the provided middlewares.
// The middlewares are applied in the order they are passed in, meaning
// the first middleware will wrap the handler last.
//
// Example usage:
//
//	http.Handle("/",
//		WrapFunc(myHandler, loggingMiddleware, authMiddleware),
//	)
func WrapFunc(h http.HandlerFunc, middlewares ...Middleware) http.Handler {
	return Wrap(h, middlewares...)
}

// Wrap wraps an http.Handler with the provided middlewares.
// Middlewares are applied in reverse order, so that the first middleware
// in the slice is executed first when handling a request.
func Wrap(h http.Handler, middlewares ...Middleware) http.Handler {
	var wrapped http.Handler = h

	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}

	return wrapped
}
