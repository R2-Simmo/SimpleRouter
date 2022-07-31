package Router

import (
	"net/http"
	"strings"
)

type router struct {
	routes    map[string]map[string]http.HandlerFunc
	options   map[string]OptionsHandler
	subRouter map[string]IRouter
	option    *Options
}

func (r router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if r.routes == nil || r.subRouter == nil { //always not hit
		panic("Router not initialized")
	}
	path := req.URL.Path
	paths := strings.Split(path, "/")
	paths = append(paths[:0], paths[1:]...)
	if !r.Exec(paths, res, req) {
		r.option.Handler.NotFound(res, req)
	}
}
func (r router) HandlerFunc(methods []string, path string, handler http.HandlerFunc) {
	paths := strings.Split(path, "/")
	paths = append(paths[:0], paths[1:]...)
	if len(paths) == 1 {
		_, exist := r.routes[paths[0]]
		if !exist {
			r.routes[paths[0]] = make(map[string]http.HandlerFunc)
		}
		for _, method := range methods {
			r.routes[paths[0]][method] = handler
		}
	} else {
		subRouter, exist := r.subRouter[paths[0]]
		if !exist { //检查子级路由是否存在
			subRouter = CreateRouter(r.option) //创建子级路由
			r.subRouter[paths[0]] = subRouter
		}
		paths = append(paths[:0], paths[1:]...)
		subRouter.HandlerFunc(methods, "/"+strings.Join(paths, "/"), handler)
	}
}
func (r router) Mount(path string, router IRouter) {
	paths := strings.Split(path, "/") //多级路由分割
	paths = append(paths[:0], paths[1:]...)
	if len(paths) == 1 {
		r.subRouter[paths[0]] = router //单一层级
	} else {
		subRouter, exist := r.subRouter[paths[0]]
		if !exist { //检查子级路由是否存在
			subRouter = CreateRouter(r.option) //创建子级路由
			r.subRouter[paths[0]] = subRouter
		}
		paths = append(paths[:0], paths[1:]...)
		subRouter.Mount("/"+strings.Join(paths, "/"), subRouter)
	}
}

func (r router) Exec(path []string, res http.ResponseWriter, req *http.Request) bool {
	if len(path) == 1 { //本级路由
		handlers, exist := r.routes[path[0]] //寻找映射的处理函数表
		if exist {
			handler, flag := handlers[req.Method] //映射到实现方法
			if flag {
				handler(res, req)
			} else {
				keys := make([]string, 0, len(handlers))
				for k := range handlers {
					keys = append(keys, k)
				}
				if req.Method == "OPTIONS" {
					options, ok := r.options[path[0]]
					if !ok {
						options = r.option.Handler.Options
					}
					options(keys, res, req)
					return true
				}
				r.option.Handler.MethodNotAllow(keys, res, req)
			}
		}
		return exist
	} else { //非本级路由
		route, exist := r.subRouter[path[0]]
		paths := append(path[:0], path[1:]...)
		if exist { //静态路径
			return route.Exec(paths, res, req) //转发次级路由执行
		}
	}
	return false
}

func (r router) GET(path string, handler http.HandlerFunc) {
	r.HandlerFunc([]string{"GET"}, path, handler)
}

func (r router) HEAD(path string, handler http.HandlerFunc) {
	r.HandlerFunc([]string{"HEAD"}, path, handler)
}

func (r router) POST(path string, handler http.HandlerFunc) {
	r.HandlerFunc([]string{"POST"}, path, handler)
}

func (r router) PUT(path string, handler http.HandlerFunc) {
	r.HandlerFunc([]string{"PUT"}, path, handler)
}

func (r router) PATCH(path string, handler http.HandlerFunc) {
	r.HandlerFunc([]string{"PATCH"}, path, handler)
}

func (r router) DELETE(path string, handler http.HandlerFunc) {
	r.HandlerFunc([]string{"DELETE"}, path, handler)
}

func (r router) OPTIONS(path string, handler OptionsHandler) {
	paths := strings.Split(path, "/")
	paths = append(paths[:0], paths[1:]...)
	if len(paths) == 1 {
		r.options[paths[0]] = handler
	} else {
		subRouter, exist := r.subRouter[paths[0]]
		if !exist { //检查子级路由是否存在
			subRouter = CreateRouter(r.option) //创建子级路由
			r.subRouter[paths[0]] = subRouter
		}
		paths = append(paths[:0], paths[1:]...)
		subRouter.OPTIONS("/"+strings.Join(paths, "/"), handler)
	}
}

func (r router) ALL(path string, handler http.HandlerFunc) {
	r.HandlerFunc([]string{"GET", "POST", "HEAD", "PUT", "PATCH", "DELETE"}, path, handler)
}

func CreateRouter(option *Options) IRouter {
	r := router{}
	r.routes = make(map[string]map[string]http.HandlerFunc)
	r.subRouter = make(map[string]IRouter)
	r.options = make(map[string]OptionsHandler)
	r.option = &Options{Handler: DefaultHandler{}}
	if option != nil {
		if option.Handler != nil {
			r.option.Handler = option.Handler
		}
	}
	return &r
}
