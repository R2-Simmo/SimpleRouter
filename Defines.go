package SimpleRouter

import (
	"net/http"
)

type IRouter interface {
	http.Handler                                      // Extends from http.Handler
	Use(module MiddlewareFunc)                        // Use sub module global
	Mount(url string, router IRouter)                 // Use sub module
	Get(url string, handlerFunc http.HandlerFunc)     // Get Method
	Post(url string, handlerFunc http.HandlerFunc)    // Post Method
	Put(url string, handlerFunc http.HandlerFunc)     // Put Method
	Delete(url string, handlerFunc http.HandlerFunc)  // Delete Method
	Patch(url string, handlerFunc http.HandlerFunc)   // Patch Method
	Head(url string, handlerFunc http.HandlerFunc)    // Head Method
	Options(url string, handlerFunc http.HandlerFunc) // Options Method
	All(url string, handlerFunc http.HandlerFunc)     // Handle All methods
}

type IHttpErrorHandler interface {
	MethodNotAllowed(w http.ResponseWriter, r *http.Request) // 405
}

type MiddlewareFunc = func(w http.ResponseWriter, r *http.Request, next NextMiddlewareFunc)

type NextMiddlewareFunc = http.HandlerFunc

type ResponseWriterWrapper struct {
	http.ResponseWriter
	completed bool
}

func (r *ResponseWriterWrapper) Write(bytes []byte) (int, error) {
	r.completed = true
	return r.ResponseWriter.Write(bytes)
}

func (r *ResponseWriterWrapper) IsCompleted() bool {
	return r.completed
}
