package main
import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"io/ioutil"
	"log"
	"io"
	"html/template"
	"os"
)

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//vid := p.ByName("vid-id")
	//vl := VIDEO_DIR + vid
	//
	//video, err := os.Open(vl)
	//defer video.Close()
	//
	//if err != nil {
	//	log.Printf("error when try to open file: %v", err)
	//	sendErrorResonse(w, http.StatusInternalServerError, "Internal ErrorChan.. 打不开视频")
	//	return
	//}
	//
	////假如文件没有扩展名，也可以把文件强制设置为mp4类型
	////w.Header().Set("Content-Type", "video/mp4")
	//http.ServeContent(w, r, "", time.Now(), video)

	//用云存储后，只用保留这些代码就行了
	log.Println("进入了streamHandler")
	targetUrl := "http://rush-videos2.oss-cn-shenzhen.aliyuncs.com/videos" + p.ByName("vid-id")
	http.Redirect(w, r, targetUrl, 301)
}


func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("/Users/zhangweicheng/Go_WS/src/video_server/upload.html")

	t.Execute(w, nil)
}



func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		log.Printf("ErrorChan，文件太大")
		msg := fmt.Sprintf("文件太大了，只接受小于%d字节的文件", MAX_UPLOAD_SIZE)
		sendErrorResonse(w, http.StatusBadRequest, msg)
		return
	}

	file, _, err := r.FormFile("file")  //<form name="file" ... />
	if err != nil {
		log.Printf("ErrorChan，文件太大，Form File解析出错")
		sendErrorResonse(w, http.StatusInternalServerError, "Internal ErrorChan")
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file err: %v。 读取上传文件出错", err)
		sendErrorResonse(w, http.StatusInternalServerError, "Internal err")
	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR + fn, data, 0666)

	if err != nil {
		log.Printf("Write file error: %v", err)
		sendErrorResonse(w, http.StatusInternalServerError, "Internal ErrorChan")
		return
	}

	ossfn := "videos/" + fn
	path := VIDEO_DIR + fn
	bn := "rush-videos2"
	ret := UploadToOss(ossfn, path, bn)
	if !ret {
		sendErrorResonse(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	os.Remove(path)

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded Successfully")
}