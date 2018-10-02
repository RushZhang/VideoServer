package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", homeHandler)

	router.POST("/", homeHandler)

	router.GET("/userhome", userHomeHandler)

	router.POST("/userhome", userHomeHandler)

	router.POST("/api", apiHandler)  //对于api的请求，做了归一化处理，代理

	router.GET("/videos/:vid-id", proxyVideoHandler)

	/*
		proxy用另一种方法解决了跨域问题，不像api归一化
		http://127.0.0.1:8080/upload/:vid-id
		http://127.0.0.1:9000/upload/:vid-id
 	*/
	router.POST("/upload/:vid-id", proxyUploadHandler)

	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))

	return router
}

func main() {
	r := RegisterHandler()
	http.ListenAndServe(":8080", r)
}



