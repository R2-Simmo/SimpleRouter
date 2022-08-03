package Router

import (
	"net/http"
	"strings"
)

type IRouter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)           //Inherit from http.Handler
	HandlerFunc([]string, string, http.HandlerFunc)         //Original request route method
	Mount(string, IRouter)                                  //Mount sub router
	Exec([]string, http.ResponseWriter, *http.Request) bool //Internal method
	GET(string, http.HandlerFunc)                           //Handle GET request
	POST(string, http.HandlerFunc)                          //Handle POST request
	PUT(string, http.HandlerFunc)                           //Handle PUT request
	DELETE(string, http.HandlerFunc)                        //Handle DELETE request
	PATCH(string, http.HandlerFunc)                         //Handle PATCH request
	HEAD(string, http.HandlerFunc)                          //Handle HEAD request
	OPTIONS(string, OptionsHandler)                         //Handle OPTIONS request
	ALL(string, http.HandlerFunc)                           //Handle all request method(exclude OPTIONS)
}
type Options struct {
	Handler EventHandler
}
type EventHandler interface {
	Forbidden(http.ResponseWriter, *http.Request)                //403
	NotFound(http.ResponseWriter, *http.Request)                 //404
	MethodNotAllow([]string, http.ResponseWriter, *http.Request) //405
	InternalError(http.ResponseWriter, *http.Request)            //500
	Options([]string, http.ResponseWriter, *http.Request)        //OPTIONS Request
}
type DefaultHandler struct {
}

func (DefaultHandler) Forbidden(res http.ResponseWriter, req *http.Request) {
	http.Error(res, http.StatusText(403), 403)
}
func (DefaultHandler) NotFound(res http.ResponseWriter, req *http.Request) {
	http.Error(res, http.StatusText(404), 404)
}
func (DefaultHandler) MethodNotAllow(allows []string, res http.ResponseWriter, req *http.Request) {
	allows = append(allows, "OPTIONS")
	res.Header().Set("Allow", strings.Join(allows, ", "))
	http.Error(res, http.StatusText(405), 405)
}
func (DefaultHandler) InternalError(res http.ResponseWriter, req *http.Request) {
	http.Error(res, http.StatusText(500), 500)
}
func (DefaultHandler) Options(allows []string, res http.ResponseWriter, req *http.Request) {
	allows = append(allows, "OPTIONS")
	res.Header().Set("Allow", strings.Join(allows, ", "))
}

type OptionsHandler func([]string, http.ResponseWriter, *http.Request)
