package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)


type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, bucketSize int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(bucketSize)
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !m.l.GetConn() { //如果没有拿到token
		sendErrorResonse(w, http.StatusTooManyRequests, "太多请求了")
		return
	}



	//好了，现在劫持完了，再做普通的handler
	m.r.ServeHTTP(w, r)

	defer m.l.ReleaseConn()
}



func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/videos/:vid-id", streamHandler)
	router.POST("/upload/:vid-id", uploadHandler)
	router.GET("/testpage", testPageHandler)

	return router
}



func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 50)
	http.ListenAndServe(":9998", mh)
}