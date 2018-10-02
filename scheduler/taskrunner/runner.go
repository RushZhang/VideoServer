package taskrunner

import (
	"fmt"
	"log"
)

/*
	runner中startDispatcher是常驻任务，一直等待runner的channel
	channel分为control channel和data channel，连接dispatcher和executor
	control chan是dispatcher和executor之间交流信息用来提醒对方
	data chan是真正用来传数据的chan
 */

type Runner struct {
	ControlChan controlChan
	ErrorChan   controlChan //Error专门用来返回CLOSE
	DataChan    dataChan
	dataSize    int
	longlived   bool  //判断是否是长期存在的Runner，如果是，则结束后不回收Runner里的字段数据
	Dispatcher  fn
	Executor    fn
}

//构造函数
func NewRunner (size int, longlived bool, d fn, e fn) *Runner {
	return &Runner {
		/*
			make(chan int,1)和make(chan int)是不一样的
			第一种通道内写入两个数据会阻塞，第二种写入就会阻塞
			如果协程在阻塞，但是主程已经退出执行，则认为程序死锁
		 */
		ControlChan: make(chan string, 1),
		ErrorChan:   make(chan string, 1),
		DataChan:    make(chan interface{}, size),
		longlived:   longlived,
		dataSize:    size,
		Dispatcher:  d,
		Executor:    e,
	}
}


func (r *Runner) startDispatch() {

	defer func() {
		if !r.longlived {
			close(r.ControlChan)
			close(r.DataChan)
			close(r.ErrorChan)
		}
	}()

	for {
		select {
			case c := <- r.ControlChan:
				if c == READY_TO_DISPATCH {
					//这一步就是进行数据分发
					err := r.Dispatcher(r.DataChan)
					if err != nil {
						r.ErrorChan <- CLOSE
					} else {
						fmt.Println("告诉executor执行")
						r.ControlChan <- READY_TO_EXECUTE
					}
				}

				if c == READY_TO_EXECUTE {
					//这一步就是开始执行数据
					err := r.Executor(r.DataChan)
					if err != nil {
						r.ErrorChan <- CLOSE
					} else {
						fmt.Println("告诉dispatcher分发")
						r.ControlChan <- READY_TO_DISPATCH
					}
				}

			case e := <- r.ErrorChan:
				if e == CLOSE {
					//这里return就是结束了startDispatch函数
					return
				}

			default:
				fmt.Println("default")

		}
	}
}



func (r *Runner) StartAll() {
	log.Println("启动了一次startAll")
	r.ControlChan <- READY_TO_DISPATCH
	r.startDispatch()
}