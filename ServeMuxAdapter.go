package SimpleRouter

import (
	"net/http"
	"strings"
)

type ServeMuxAdapter struct {
	serveMux          *http.ServeMux
	errorHandler      IHttpErrorHandler
	routes            map[string][]IRouter
	middlewareHandler *MiddlewareHandler
}

func (h *ServeMuxAdapter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.middlewareHandler.HandleMiddleware(writer, request, func(w http.ResponseWriter, r *http.Request) {
		// 子路由处理
		for key := range h.routes {
			if strings.HasPrefix(r.URL.Path, key) {
				for _, subRouter := range h.routes[key] {
					subRouter.ServeHTTP(w.(*ResponseWriterWrapper), r)
					if w.(*ResponseWriterWrapper).IsCompleted() {
						return
					}
				}
			}
		}
		h.serveMux.ServeHTTP(w, r)
	})
}

func (h *ServeMuxAdapter) handle(s string, handlerFunc http.HandlerFunc) {
	h.serveMux.HandleFunc(s, handlerFunc)
}

func (h *ServeMuxAdapter) Use(module MiddlewareFunc) {
	h.middlewareHandler.RegisterMiddleware(module)
}

func (h *ServeMuxAdapter) Mount(url string, router IRouter) {
	if url == "" {
		url = "/" // mount to root scope
	}
	h.routes[url] = append(h.routes[url], router)
}

func (h *ServeMuxAdapter) Get(url string, handlerFunc http.HandlerFunc) {
	h.serveMux.HandleFunc(url, h.handlerWrap("get", handlerFunc))
}

func (h *ServeMuxAdapter) Post(url string, handlerFunc http.HandlerFunc) {
	h.handle(url, h.handlerWrap("post", handlerFunc))
}

func (h *ServeMuxAdapter) Put(url string, handlerFunc http.HandlerFunc) {
	h.handle(url, h.handlerWrap("PUT", handlerFunc))
}

func (h *ServeMuxAdapter) Delete(url string, handlerFunc http.HandlerFunc) {
	h.handle(url, h.handlerWrap("DELETE", handlerFunc))
}

func (h *ServeMuxAdapter) Patch(url string, handlerFunc http.HandlerFunc) {
	h.handle(url, h.handlerWrap("PATCH", handlerFunc))
}

func (h *ServeMuxAdapter) Head(url string, handlerFunc http.HandlerFunc) {
	h.handle(url, h.handlerWrap("HEAD", handlerFunc))
}

func (h *ServeMuxAdapter) Options(url string, handlerFunc http.HandlerFunc) {
	h.handle(url, h.handlerWrap("OPTIONS", handlerFunc))
}

func (h *ServeMuxAdapter) All(url string, handlerFunc http.HandlerFunc) {
	h.handle(url, handlerFunc)
}

func (h *ServeMuxAdapter) handlerWrap(method string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO merge handler in same url
		if strings.ToUpper(r.Method) != strings.ToUpper(method) {
			h.errorHandler.MethodNotAllowed(w, r)
			return
		}
		handlerFunc(w, r)
	}
}

func NewServeMuxAdapter() IRouter {
	router := &ServeMuxAdapter{serveMux: http.NewServeMux(), errorHandler: DefaultHttpErrorHandler{}, routes: make(map[string][]IRouter), middlewareHandler: NewMiddlewareHandler()}
	return router
}
