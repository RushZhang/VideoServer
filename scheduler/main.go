package main

import (
	"github.com/julienschmidt/httprouter"
	"video_server/scheduler/taskrunner"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video-delete-record/:vid-id", vidDelRecHandler)

	return router
}

func main() {
	//c := make(chan int)
	go taskrunner.Start()
	r := RegisterHandlers()
	//<- c
	http.ListenAndServe(":9001", r)  //因为http.ListenAndServe是阻塞的，如果没有这个，就需要添加上面两行注释的代码
}
