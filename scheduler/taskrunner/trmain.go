package taskrunner

import (
	"time"
)

/*
	timer
	setup
	start{trigger -> task -> runner}
 */


type Worker struct {
	ticker *time.Ticker
	runner *Runner
}


func NewWorker(interval time.Duration, r *Runner) *Worker {

	return &Worker{
		ticker: time.NewTicker(interval *time.Second),
		runner: r,
	}

}

func (w *Worker) startWorker() {

	for {
		select {
			case <- w.ticker.C:   //ticker每过interval的时间就会发送一个这种消息
				go w.runner.StartAll()

		}
	}

	//不能用下面的方法，因为range会有延迟，时间久了定时器就不准了
	//for c = range w.ticker.C {
	//
	//}
}


func Start() {
	//start video file cleaning
	r := NewRunner(3, true, VideoClearDispatcher, VideoClearExecutor)
	w := NewWorker(100, r)
	go w.startWorker()
}