# SimpleRouter

简体中文/[English](README.md)

这是一个简单的用于GO语言Web服务的路由分发组件

## 模块功能

### 路由中间件

中间件系统的行为与栈接近,
系统会默认优先调用注册较晚的中间件,
并参考洋葱模型对请求上下文进行传播,
最终传递至路由系统.

简单来说,
可以使用 `IRouter.Use` 注册中间件,
参考下面这个例子:

```go
// ...

func MyMiddleware(w http.ResponseWriter, r *http.Request, next SimpleRouter.NextMiddlewareFunc) {
	// 路由预处理...
	next(w, r)
	// 路由后处理...
}

// ...

func main() {
	router := SimpleRouter.NewServeMuxAdapter()
	router.Use(MyMiddleware)
	// ...
	// 业务逻辑和其他路由
	// ...
	http.ListenAndServe("localhost:8080", router)
	// ...
}
```

完整的使用示例可在 `examples/Middleware.go` 找到,
该文件中使用了封装的 `http.ServeMux` 作为基本的路由功能实现.