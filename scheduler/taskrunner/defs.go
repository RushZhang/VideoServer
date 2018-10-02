package taskrunner


const (
	READY_TO_DISPATCH = "d"
	READY_TO_EXECUTE = "e"
	CLOSE = "c"
	VIDEO_PATH = "/Users/zhangweicheng/Go_WS/src/video_server/videos/"
)

type controlChan chan string

type dataChan chan interface{}

//fn就是我们的dispatcher或executor
type fn func(dc dataChan) error