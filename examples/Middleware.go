package main

import (
	"SimpleRouter"
	"log"
	"net/http"
	"strconv"
	"time"
)

func RequestTimeTrace(w http.ResponseWriter, r *http.Request, next SimpleRouter.NextMiddlewareFunc) {
	currentTime := time.Now()
	log.Printf("[%s] %s %s\n", currentTime.Format("2006-01-02 15:04:05"), r.Method, r.RequestURI)
	next(w, r)
}
func ResponseTimeTrace(w http.ResponseWriter, r *http.Request, next SimpleRouter.NextMiddlewareFunc) {
	currentTime := time.Now()
	next(w, r)
	responseTime := time.Now().Sub(currentTime).Milliseconds()
	log.Printf("Response Time: %dms", responseTime)
	// No effect cause Go HTTP policy
	w.Header().Add("ResponseTime", strconv.FormatInt(responseTime, 10))
}
func main() {
	router := SimpleRouter.NewServeMuxAdapter()
	router.Use(RequestTimeTrace)
	router.Use(ResponseTimeTrace)
	router.All("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("Hello, World!"))
		if err != nil {
			return
		}
	})
	http.ListenAndServe("localhost:8080", router)
}
