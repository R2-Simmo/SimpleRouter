package Router

import (
	"net/http"
	"strings"
)

type IRouter interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	HandlerFunc([]string, string, http.HandlerFunc)
	Mount(string, IRouter)
	Exec([]string, http.ResponseWriter, *http.Request) bool
	GET(string, http.HandlerFunc)
	POST(string, http.HandlerFunc)
	PUT(string, http.HandlerFunc)
	DELETE(string, http.HandlerFunc)
	PATCH(string, http.HandlerFunc)
	HEAD(string, http.HandlerFunc)
	OPTIONS(string, OptionsHandler)
	ALL(string, http.HandlerFunc)
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
