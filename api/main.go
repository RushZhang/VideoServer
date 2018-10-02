package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"video_server/api/session"
)



/*
	httprouter.Router是实现了Handler这个接口
	Handler里有个HttpServe函数
	所以我们在这里做点手脚就可以提前处理http请求
 */
type middleWareHandler struct {
	r *httprouter.Router
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.r.ServeHTTP(w, r)
}



//这个函数有点像是构造器的感觉
func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}




func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)

	router.POST("/user/:username", Login)

	router.GET("/user/:username", GetUserInfo)

	router.POST("/user/:username/videos", AddNewVideo)

	router.GET("/user/:username/videos", ListAllVideos)

	router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)

	router.POST("/videos/:vid-id/comments", PostComment)

	router.GET("/videos/:vid-id/comments", ShowComments)

	return router
}


func Prepare() {
	session.LoadSessionsFromDB()
}


func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	//这里handler已经是被middleware做过手脚的了
	fmt.Println("正在8000端口监听。。。", mh)
	http.ListenAndServe(":8000", r)

}


/*
	基本流程：
	main -> middleware(校验、流控...) -> defs(包括msg, err) -> handlers -> dbops -> response
 */

