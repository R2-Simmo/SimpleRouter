package SimpleRouter

import (
	"net/http"
	"sync"
)

type MiddlewareHandler struct {
	middlewares []MiddlewareFunc
}

func NewMiddlewareHandler() *MiddlewareHandler {
	return &MiddlewareHandler{middlewares: make([]MiddlewareFunc, 0)}
}

func (h *MiddlewareHandler) HandleMiddleware(writer http.ResponseWriter, request *http.Request, fn func(w http.ResponseWriter, r *http.Request)) {
	// 中间件调用
	next := make([]NextMiddlewareFunc, len(h.middlewares)+1)
	next[0] = func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
		return
	}
	mutex := make([]*sync.Mutex, len(h.middlewares))
	complete := sync.Mutex{}
	for i := 0; i < len(mutex); i++ {
		m := &sync.Mutex{}
		mutex[i] = m
	}
	complete.Lock()
	for i, middleware := range h.middlewares {
		index := i
		m := middleware
		next[i+1] = func(w http.ResponseWriter, r *http.Request) {
			mutex[index].Lock()
			m(w, r, next[index])
			mutex[index].Unlock()
			if index == len(mutex)-1 {
				complete.Unlock()
			}
		}
	}
	go next[len(h.middlewares)](writer, request)
	complete.Lock()
	complete.Unlock()
}

func (h *MiddlewareHandler) RegisterMiddleware(middleware MiddlewareFunc) {
	h.middlewares = append(h.middlewares, middleware)
}
