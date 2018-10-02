package taskrunner

import (
	"testing"
	"log"
	"time"
)

func TestRunner(t *testing.T) {
	d := func(dc dataChan) error {
		for i := 0; i < 3; i++ {
			dc <- i
			log.Printf("Dispatch发送了: %v", i)
		}
		return nil
	}

	e := func(dc dataChan) error {
		forloop:
			for {
				select {
					case d := <- dc:
						log.Printf("Executor得到了: %v", d)
					//当得不到数据时候就会进入default
					default:
						break forloop
				}
			}
		return nil
		//return errors.New("结束咯，测试关闭dispatcher")
	}


	runner := NewRunner(30, false, d, e)

	//因为StartDispatcher里有无限循环，所以只能用go启动
	go runner.StartAll()

	time.Sleep(3 * time.Second)
}
