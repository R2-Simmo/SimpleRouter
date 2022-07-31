package Router

import (
	"net/http"
)

type IRouter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Mount(string, IRouter)
	Exec([]string, http.ResponseWriter, *http.Request) bool
	Get(string, func(http.ResponseWriter, *http.Request))
	Head(string, func(http.ResponseWriter, *http.Request))
	Post(string, func(http.ResponseWriter, *http.Request))
	Put(string, func(http.ResponseWriter, *http.Request))
	Patch(string, func(http.ResponseWriter, *http.Request))
	Delete(string, func(http.ResponseWriter, *http.Request))
	Trace(string, func(http.ResponseWriter, *http.Request))
	All(string, func(http.ResponseWriter, *http.Request))
}
type Options struct {
	Handler ErrorHandler
}
type ErrorHandler interface {
	Forbidden(http.ResponseWriter, *http.Request)      //403
	NotFound(http.ResponseWriter, *http.Request)       //404
	MethodNotAllow(http.ResponseWriter, *http.Request) //405
	InternalError(http.ResponseWriter, *http.Request)  //500
}
type DefaultHandler struct {
}

func (DefaultHandler) Forbidden(res http.ResponseWriter, req *http.Request) {
	http.Error(res, http.StatusText(403), 403)
}
func (DefaultHandler) NotFound(res http.ResponseWriter, req *http.Request) {
	http.Error(res, http.StatusText(404), 404)
}
func (DefaultHandler) MethodNotAllow(res http.ResponseWriter, req *http.Request) {
	http.Error(res, http.StatusText(405), 405)
}
func (DefaultHandler) InternalError(res http.ResponseWriter, req *http.Request) {
	http.Error(res, http.StatusText(500), 500)
}
